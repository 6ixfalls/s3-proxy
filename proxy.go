package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Proxy interface {
	Get(key string) (*s3.GetObjectOutput, error)
	GetWebsiteConfig() (*s3.GetBucketWebsiteOutput, error)
}

type RealS3Proxy struct {
	bucket string
	s3     *s3.S3
}

func NewS3Proxy(key, secret, region, bucket string) S3Proxy {
	defaultResolver := endpoints.DefaultResolver()
	s3CustResolverFn := func(service, region string, optFns ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
		if service == "s3" {
			return endpoints.ResolvedEndpoint{
				URL:           "https://s3.filebase.com",
				SigningRegion: "us-east-1",
			}, nil
		}

		return defaultResolver.EndpointFor(service, region, optFns...)
	}
	
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(key, secret, ""),
		EndpointResolver: endpoints.ResolverFunc(s3CustResolverFn),
	}))

	return &RealS3Proxy{
		bucket: bucket,
		s3:     s3.New(sess),
	}
}

func (p *RealS3Proxy) Get(key string) (*s3.GetObjectOutput, error) {
	req := &s3.GetObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(key),
	}

	return p.s3.GetObject(req)
}

func (p *RealS3Proxy) GetWebsiteConfig() (*s3.GetBucketWebsiteOutput, error) {
	req := &s3.GetBucketWebsiteInput{
		Bucket: aws.String(p.bucket),
	}

	return p.s3.GetBucketWebsite(req)
}
