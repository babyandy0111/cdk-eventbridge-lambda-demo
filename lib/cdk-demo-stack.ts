import * as cdk from '@aws-cdk/core';
import * as iam from '@aws-cdk/aws-iam';
import * as events from '@aws-cdk/aws-events';
import * as logs from '@aws-cdk/aws-logs';
import * as targets from '@aws-cdk/aws-events-targets';
import * as path from 'path';
import * as apigw from '@aws-cdk/aws-apigateway';
import * as sns from '@aws-cdk/aws-sns';
import * as subscriptions from "@aws-cdk/aws-sns-subscriptions";
import * as cw_actions from "@aws-cdk/aws-cloudwatch-actions";
import * as ec2 from '@aws-cdk/aws-ec2';
import * as cloudwatch from '@aws-cdk/aws-cloudwatch';

export class CdkDemoStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const lambda = require('@aws-cdk/aws-lambda');
    const stageName = ''
    const webhookEndpointURL = ''
    const mayoEndpointURL = ''
    const fs_domain = ''
    const fsKey = ''
    const mayo_api_domain = ''
    const mayo_key = ''
    const TENANT_ID = ''
    const CLIENT_ID = ''
    const CLIENT_SECRET = ''
    const AAD_ENDPOINT = ''
    const GRAPH_ENDPOINT = ''

    // 這邊是給監控寶使用的lambda
    const webhookLambda = new lambda.Function(this, 'webhook', {
      runtime: lambda.Runtime.GO_1_X,
      timeout: cdk.Duration.seconds(300),
      handler: "main",
      code: lambda.Code.fromAsset(
        path.join(__dirname, "../webhook")
      ),
      environment: {
        region: cdk.Stack.of(this).region,
        zones: JSON.stringify(cdk.Stack.of(this).availabilityZones),
        fs_domain: fs_domain,
        fsKey: fsKey,
      },
    })

    const restApi = new apigw.RestApi(this, "dev-api", {deploy: false});

    restApi.root
      .addResource(webhookEndpointURL)
      .addMethod("GET", new apigw.LambdaIntegration(webhookLambda, {proxy: true}));

    const devDeploy = new apigw.Deployment(this, "dev-deployment", {api: restApi});

    // 定義stage => 網址會是 /dev/bot
    new apigw.Stage(this, "msp-stage", {
      deployment: devDeploy,
      stageName: stageName
    });

    new cdk.CfnOutput(this, "webhook URL", {
      value: `https://${restApi.restApiId}.execute-api.${this.region}.amazonaws.com/${stageName}/${webhookEndpointURL}`,
    });

    // 這邊是mayo api
    const mayoLambda = new lambda.Function(this, 'mayo', {
      runtime: lambda.Runtime.NODEJS_16_X,
      timeout: cdk.Duration.seconds(300),
      handler: "index.handler",
      code: lambda.Code.fromAsset(
        path.join(__dirname, "../mayo")
      ),
      environment: {
        region: cdk.Stack.of(this).region,
        zones: JSON.stringify(cdk.Stack.of(this).availabilityZones),
        mayo_api_domain: mayo_api_domain,
        mayo_key: mayo_key,
        tenant_id: TENANT_ID,
        client_id: CLIENT_ID,
        client_secret: CLIENT_SECRET,
        aad_endpoint: AAD_ENDPOINT,
        graph_endpoint: GRAPH_ENDPOINT,
      }
    })

    const mayoApi = new apigw.RestApi(this, "dev-mayo-api", {deploy: false});
    mayoApi.root
      .addResource(mayoEndpointURL)
      .addMethod("GET", new apigw.LambdaIntegration(mayoLambda, {proxy: true}));

    const devMayoDeploy = new apigw.Deployment(this, "dev-mayo-deployment", {api: mayoApi});

    new apigw.Stage(this, "mayo-stage", {
      deployment: devMayoDeploy,
      stageName: stageName
    });

    new cdk.CfnOutput(this, "mayo URL", {
      value: `https://${mayoApi.restApiId}.execute-api.${this.region}.amazonaws.com/${stageName}/${mayoEndpointURL}`,
    });

    // freshService
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
        fs_domain: fs_domain,
        fsKey: fsKey,
      },
    })

    // iam
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
        fs_domain: fs_domain,
        fsKey: fsKey,
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
    const whLogGroup = new logs.LogGroup(this, 'whLogGroup', {
      logGroupName: 'whLogGroup',
    })

    const fsLogGroup = new logs.LogGroup(this, 'fsLogGroup', {
      logGroupName: 'fsLogGroup',
    })

    const iamLogGroup = new logs.LogGroup(this, 'iamLogGroup', {
      logGroupName: 'iamLogGroup',
    })

    // add fs eventBridge
    const mspEventRule = new events.Rule(this, 'mspRule', {
      eventPattern: {
        source: ["aws.support"],
        detailType: ["Support Case Update"],
      },
    })

    mspEventRule.addTarget(new targets.CloudWatchLogGroup(fsLogGroup))
    mspEventRule.addTarget(new targets.LambdaFunction(fsLambda))

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

    // for test 用
    const cfnKeyPair = new ec2.CfnKeyPair(this, 'mspCfnKeyPair', {
      keyName: 'KeyPair',
    });

    const vpc = ec2.Vpc.fromLookup(this, "VPC", {
      isDefault: true,
    });

    const securityGroup = new ec2.SecurityGroup(this, "SecurityGroup", {
      vpc,
      allowAllOutbound: true,
    });
    securityGroup.addIngressRule(ec2.Peer.anyIpv4(), ec2.Port.tcp(22));

    const instance = new ec2.Instance(this, "Instance", {
      vpc,
      instanceType: ec2.InstanceType.of(
        ec2.InstanceClass.T3A,
        ec2.InstanceSize.NANO
      ),
      machineImage: ec2.MachineImage.latestAmazonLinux({
        generation: ec2.AmazonLinuxGeneration.AMAZON_LINUX_2,
      }),
      securityGroup,
      vpcSubnets: {
        subnetType: ec2.SubnetType.PUBLIC,
      },
      keyName: "KeyPair",
    });

    // 增加一個topic
    const topic = new sns.Topic(this, "mspTopic");
    const metric = new cloudwatch.Metric({
      namespace: "AWS/EC2",
      metricName: "CPUUtilization",
      dimensions: {
        InstanceId: instance.instanceId,
      },
      period: cdk.Duration.minutes(1),
    });

    const alarm = new cloudwatch.Alarm(this, "Alarm", {
      metric,
      threshold: 5,
      evaluationPeriods: 1,
    });

    alarm.addAlarmAction(new cw_actions.SnsAction(topic));

    const alarmLambda = new lambda.Function(this, 'alarm', {
      runtime: lambda.Runtime.GO_1_X,
      timeout: cdk.Duration.seconds(300),
      handler: "main",
      code: lambda.Code.fromAsset(
        path.join(__dirname, "../alarm")
      ),
      environment: {
        region: cdk.Stack.of(this).region,
        zones: JSON.stringify(cdk.Stack.of(this).availabilityZones),
        fs_domain: fs_domain,
        fsKey: fsKey,
      },
    })

    topic.addSubscription(new subscriptions.LambdaSubscription(alarmLambda));

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
