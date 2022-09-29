import * as cdk from '@aws-cdk/core';
// import * as lambda from '@aws-cdk/aws-lambda';
import * as iam from '@aws-cdk/aws-iam';
import * as events from '@aws-cdk/aws-events';
import * as logs from '@aws-cdk/aws-logs';
import * as targets from '@aws-cdk/aws-events-targets';
import * as path from 'path';
// import * as sns from '@aws-cdk/aws-sns';
// import * as subscriptions from "@aws-cdk/aws-sns-subscriptions";


export class CdkDemoStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const lambda = require('@aws-cdk/aws-lambda');

    const fsLambda = new lambda.Function(this, 'fs', {
      runtime: lambda.Runtime.GO_1_X,
      timeout: cdk.Duration.seconds(300),
      handler: "main",
      code: lambda.Code.fromAsset(
        path.join(__dirname, "../fs")
      ),
      environment: {
        region: cdk.Stack.of(this).region,
        zones: JSON.stringify(cdk.Stack.of(this).availabilityZones),
        fs_domain: process.env.FS_DOMAIN,
        fsKey: process.env.FSKEY,
      },
    })

    const iamLambda = new lambda.Function(this, 'check-iam', {
      runtime: lambda.Runtime.GO_1_X,
      timeout: cdk.Duration.seconds(300),
      handler: "main",
      code: lambda.Code.fromAsset(
        path.join(__dirname, "../iam")
      ),
      environment: {
        region: cdk.Stack.of(this).region,
        zones: JSON.stringify(cdk.Stack.of(this).availabilityZones),
        fs_domain: process.env.FS_DOMAIN,
        fsKey: process.env.FSKEY,
      },
    })

    // add policy
    const supportPolicy = new iam.PolicyStatement({
      actions: ['support:*'],
      resources: ['*'],
    })

    fsLambda.role?.attachInlinePolicy(
      new iam.Policy(this, 'fsSupportPolicy', {
        statements: [supportPolicy],
      }),
    )

    iamLambda.role?.attachInlinePolicy(
      new iam.Policy(this, 'iamSupportPolicy', {
        statements: [supportPolicy],
      }),
    )

    // cloudwatch log group
    const fsLogGroup = new logs.LogGroup(this, 'fsLogGroup', {
      logGroupName: 'fsLogGroup',
    })

    const iamLogGroup = new logs.LogGroup(this, 'iamLogGroup', {
      logGroupName: 'iamLogGroup',
    })

    // add eventBridge
    const supportEventRule = new events.Rule(this, 'supportRule', {
      eventPattern: {
        source: ["aws.support"],
        detailType: ["Support Case Update"],
      },
    })

    supportEventRule.addTarget(new targets.CloudWatchLogGroup(fsLogGroup))
    supportEventRule.addTarget(new targets.LambdaFunction(fsLambda))

    // add check iam eventBridge
    const iamEventRule = new events.Rule(this, 'iamRule', {
      eventPattern: {
        source: ["aws.trustedadvisor"],
        detail: {
          "status": [
            "ERROR"
          ],
          "check-name": [
            "Exposed Access Keys"
          ]
        },
        detailType: ["Trusted Advisor Check Item Refresh Notification"],
      },
    })

    iamEventRule.addTarget(new targets.CloudWatchLogGroup(iamLogGroup))
    iamEventRule.addTarget(new targets.LambdaFunction(iamLambda))

    // const topic = new sns.Topic(this, 'fsTopic', {
    //   contentBasedDeduplication: true,
    //   displayName: 'fs subscription topic',
    //   fifo: true,
    //   topicName: 'aws2fs',
    // })
    //
    // topic.addSubscription(new subscriptions.LambdaSubscription(fsLambda))
  }
}
