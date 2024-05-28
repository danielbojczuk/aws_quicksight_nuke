package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/quicksight"
)

func main() {
	region := "eu-central-1"
	accoundID := "ACCOUNT_ID"
	activeSubscriptionStatus := "ACCOUNT_CREATED"
	terminateProtectionEnabled := false


	sessionOpts := session.Options{}
	sessionOpts.Config.Region = &region


	sess, err := session.NewSessionWithOptions(sessionOpts)
	
	if(err != nil) {
		log.Fatal(err)
	}
	svc := quicksight.New(sess)

	describeSubscritpion := quicksight.DescribeAccountSubscriptionInput {
		AwsAccountId: &accoundID,
	}

	describeSubscriptionOutput, err := svc.DescribeAccountSubscription(&describeSubscritpion)
	if(err != nil) {
		var resoureceNotFoundException *quicksight.ResourceNotFoundException
		if(!errors.As(err,&resoureceNotFoundException)) {
			log.Fatal(err)
		}

		log.Println("No Quicksight subscription was found for this account")
		return
	}
	if(*describeSubscriptionOutput.AccountInfo.AccountSubscriptionStatus != activeSubscriptionStatus) {
		log.Println("There is No active Quicksight subscription for this account ")
		return
	}


	describeSettingsInput := quicksight.DescribeAccountSettingsInput {
		AwsAccountId: &accoundID,
	}

	describeSettingsOutput, err := svc.DescribeAccountSettings(&describeSettingsInput)
	if(err != nil) {
		log.Fatal(err)
	}

	if(*describeSettingsOutput.AccountSettings.TerminationProtectionEnabled) {
		log.Println("Disabling termination protection for Quicksight subscription")
		updateSettingsInput := quicksight.UpdateAccountSettingsInput{
			AwsAccountId: &accoundID,
			DefaultNamespace: describeSettingsOutput.AccountSettings.DefaultNamespace,
			NotificationEmail: describeSettingsOutput.AccountSettings.NotificationEmail,
			TerminationProtectionEnabled: &terminateProtectionEnabled,
		}

		_, err = svc.UpdateAccountSettings(&updateSettingsInput)
		if(err != nil) {
			log.Fatal(err)
		}
	}

	deleteSubscriptionInput := quicksight.DeleteAccountSubscriptionInput {
		AwsAccountId: &accoundID,
	}

	deleteSubscriptionOutput, err := svc.DeleteAccountSubscription(&deleteSubscriptionInput)

	if(err != nil) {
		log.Fatal(err)
	}

	fmt.Println(*deleteSubscriptionOutput)
}	
