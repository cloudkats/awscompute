package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	_ "github.com/aws/aws-sdk-go-v2/service/sts"
)

type Account struct {
	Account string
	Aliases string
}

func (cmp Compute) AccountAnalyzer() (*Account, error) {
	client := iam.NewFromConfig(cmp.Config)

	result, err := GetAccountAliases(cmp.Context, client, &iam.ListAccountAliasesInput{MaxItems: aws.Int32(int32(100))})
	if err != nil {
		return nil, err
	}

	aliases := ""
	for _, alias := range result.AccountAliases {
		aliases += alias
	}

	test := sts.NewFromConfig(cmp.Config)

	r, err := test.GetCallerIdentity(cmp.Context, &sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, err
	}

	return &Account{
		Aliases: aliases,
		Account: *r.Account,
	}, nil
}

// https://aws.github.io/aws-sdk-go-v2/docs/code-examples/iam/listaccountaliases/

// IAMListAccountAliasesAPI defines the interface for the ListAccountAliases function.
// We use this interface to test the function using a mocked service.
type IAMListAccountAliasesAPI interface {
	ListAccountAliases(ctx context.Context,
		params *iam.ListAccountAliasesInput,
		optFns ...func(*iam.Options)) (*iam.ListAccountAliasesOutput, error)
}

// GetAccountAliases retrieves the aliases for your AWS Identity and Access Management (IAM) account.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If successful, a ListAccountAliasesOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to ListAccountAliases.
func GetAccountAliases(c context.Context, api IAMListAccountAliasesAPI, input *iam.ListAccountAliasesInput) (*iam.ListAccountAliasesOutput, error) {
	return api.ListAccountAliases(c, input)
}
