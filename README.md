# talos-test

https://www.talos.dev/v1.6/talos-guides/install/cloud-platforms/aws

```
➜  ~ REGION=us-east-1                                                                                                        
➜  ~ aws --profile development ec2 describe-vpcs --region $REGION                                                            
{                                                                                                                            
    "Vpcs": [                                                                                                                
        {                                                                                                                    
            "CidrBlock": "10.200.0.0/16",                                                                                    
            "DhcpOptionsId": "dopt-0d8b2edc45ff5bbcf",                                                                       
            "State": "available",                                                                                            
            "VpcId": "vpc-0d2df8f32c707fd36",                                                                                
            "OwnerId": "339735964233",                                                                                       
            "InstanceTenancy": "default",                                                                                    
            "CidrBlockAssociationSet": [                                                                                     
                {                                                                                                            
                    "AssociationId": "vpc-cidr-assoc-02a5fe171cd4c3f2b",                                                     
                    "CidrBlock": "10.200.0.0/16",                                                                            
                    "CidrBlockState": {                                                                                      
                        "State": "associated"                                                                                
                    }                                                                                                        
                }                                                                                                            
            ],                                                                                                               
            "IsDefault": false,                                                                                              
            "Tags": [                                                                                                        
                {
                    "Key": "Name",
                    "Value": "development"
                }
            ]
        }
    ]
}
➜  ~ VPC
➜  ~ VPC=vpc-0d2df8f32c707fd36
```

```
➜  ~ AMI=`curl -sL https://github.com/siderolabs/talos/releases/download/v1.6.4/cloud-images.json | \
    jq -r '.[] | select(.region == "'$REGION'") | select (.arch == "amd64") | .id'`
echo $AMI
ami-09360283b6eec5d54
```

```
➜  ~ aws --profile development ec2 create-security-group \        
    --region $REGION \
    --group-name talos-aws-tutorial-sg \
    --description "Security Group for EC2 instances to allow ports required by Talos"

SECURITY_GROUP="(security group id that is returned)"                                   

An error occurred (VPCIdNotSpecified) when calling the CreateSecurityGroup operation: No default VPC for this user
```


```
➜  code aws --profile development ec2 create-security-group \
    --region $REGION \
    --group-name talos-aws-tutorial-sg \
    --description "Security Group for EC2 instances to allow ports required by Talos" \
    --vpc-id $VPC

SECURITY_GROUP="(security group id that is returned)"
{
    "GroupId": "sg-01e3966613a73846c"
}
```

```
➜  code SECURITY_GROUP=sg-01e3966613a73846c
➜  code aws --profile development ec2 authorize-security-group-ingress \
    --region $REGION \
    --group-id $SECURITY_GROUP \
    --protocol all \
    --port 0 \
    --source-group $SECURITY_GROUP
{
    "Return": true,
    "SecurityGroupRules": [
        {
            "SecurityGroupRuleId": "sgr-03e616bc41627cdc7",
            "GroupId": "sg-01e3966613a73846c",
            "GroupOwnerId": "339735964233",
            "IsEgress": false,
            "IpProtocol": "-1",
            "FromPort": -1,
            "ToPort": -1,
            "ReferencedGroupInfo": {
                "GroupId": "sg-01e3966613a73846c",
                "UserId": "339735964233"
            }
        }
    ]
}
```

```
aws --profile development ec2 authorize-security-group-ingress \
    --region $REGION \
    --group-id $SECURITY_GROUP \
    --protocol tcp \
    --port 6443 \
    --cidr 0.0.0.0/0

aws --profile development ec2 authorize-security-group-ingress \
    --region $REGION \
    --group-id $SECURITY_GROUP \
    --protocol tcp \
    --port 50000-50001 \
    --cidr 0.0.0.0/0
```

```
➜  code aws --profile development ec2 authorize-security-group-ingress \                                                                                                                   
    --region $REGION \                                                                                                                                                
    --group-id $SECURITY_GROUP \                              
    --protocol tcp \                                                
    --port 6443 \                                                          
    --cidr 0.0.0.0/0                                          
                                                                    
aws --profile development ec2 authorize-security-group-ingress \                                                                                                                           
    --region $REGION \                                                     
    --group-id $SECURITY_GROUP \                                                   
    --protocol tcp \                                                       
    --port 50000-50001 \                                                   
    --cidr 0.0.0.0/0                                                               
{                                                             
    "Return": true,                                           
    "SecurityGroupRules": [                                   
        {                                                           
            "SecurityGroupRuleId": "sgr-08e2fbc5350e5e933",   
            "GroupId": "sg-01e3966613a73846c",                      
            "GroupOwnerId": "339735964233",                                
            "IsEgress": false,                                
            "IpProtocol": "tcp",                                    
            "FromPort": 6443,                                              
            "ToPort": 6443,                                                
            "CidrIpv4": "0.0.0.0/0"                                                
        }                                                                  
    ]                                                                      
}                                                                                  
{                                                                          
    "Return": true,                  
    "SecurityGroupRules": [                                                        
        {                            
            "SecurityGroupRuleId": "sgr-0c86de4a8e13b0d15",                                  
            "GroupId": "sg-01e3966613a73846c",                                     
            "GroupOwnerId": "339735964233",                                                  
            "IsEgress": false,           
            "IpProtocol": "tcp",                                                             
            "FromPort": 50000,           
            "ToPort": 50001,                  
            "CidrIpv4": "0.0.0.0/0"                                                          
        }                                     
    ]                                         
}  
```

```
aws --profile development elbv2 create-load-balancer \
    --region $REGION \
    --name talos-aws-tutorial-lb \
    --type network --subnets $SUBNET
```

```
➜  code aws --profile development elbv2 create-load-balancer \
    --region $REGION \
    --name talos-aws-tutorial-lb \
    --type network --subnets $SUBNET
{
    "LoadBalancers": [
        {
            "LoadBalancerArn": "arn:aws:elasticloadbalancing:us-east-1:339735964233:loadbalancer/net/talos-aws-tutorial-lb/eedbb2eaa739616a",
            "DNSName": "talos-aws-tutorial-lb-eedbb2eaa739616a.elb.us-east-1.amazonaws.com",
            "CanonicalHostedZoneId": "Z26RNL4JYFTOTI",
            "CreatedTime": "2024-03-05T22:56:05.362000+00:00",
            "LoadBalancerName": "talos-aws-tutorial-lb",
            "Scheme": "internet-facing",
            "VpcId": "vpc-0d2df8f32c707fd36",
            "State": {
                "Code": "provisioning"
            },
            "Type": "network",
            "AvailabilityZones": [
                {
                    "ZoneName": "us-east-1a",
                    "SubnetId": "subnet-0bdf61de47f04f337",
                    "LoadBalancerAddresses": []
                }
            ],
            "IpAddressType": "ipv4"
        }
    ]
}
```
```
LOAD_BALANCER_ARN=arn:aws:elasticloadbalancing:us-east-1:339735964233:loadbalancer/net/talos-aws-tutorial-lb/eedbb2eaa739616a
```
```
➜  code aws --profile development elbv2 create-target-group \ 
    --region $REGION \
    --name talos-aws-tutorial-tg \
    --protocol TCP \
    --port 6443 \
    --target-type ip \
    --vpc-id $VPC                            
{
    "TargetGroups": [
        {
            "TargetGroupArn": "arn:aws:elasticloadbalancing:us-east-1:339735964233:targetgroup/talos-aws-tutorial-tg/7146ea36aa0472cc",
            "TargetGroupName": "talos-aws-tutorial-tg",
            "Protocol": "TCP",
            "Port": 6443,
            "VpcId": "vpc-0d2df8f32c707fd36",
            "HealthCheckProtocol": "TCP",
            "HealthCheckPort": "traffic-port",
            "HealthCheckEnabled": true,
            "HealthCheckIntervalSeconds": 30,
            "HealthCheckTimeoutSeconds": 10,
            "HealthyThresholdCount": 5,
            "UnhealthyThresholdCount": 2,
            "TargetType": "ip",
            "IpAddressType": "ipv4"
        }
    ]
}
```
```
TARGET_GROUP_ARN=arn:aws:elasticloadbalancing:us-east-1:339735964233:targetgroup/talos-aws-tutorial-tg/7146ea36aa0472cc
```


```
➜  code talosctl gen config talos-k8s-aws-tutorial https://talos-aws-tutorial-lb-eedbb2eaa739616a.elb.us-east-1.amazonaws.com:443 --with-examples=false --with-docs=false
generating PKI and tokens
Created /Users/samebstein/tkhq/code/controlplane.yaml
Created /Users/samebstein/tkhq/code/worker.yaml
Created /Users/samebstein/tkhq/code/talosconfig
➜  code ls
controlplane.yaml   gitops              infrastructure      mono                qos                 talosconfig         worker.yaml
docs                go-sdk              keys                mono-onboarding     talos-test          terraform-aws-talos
➜  code mv controlplane.yaml talos-test/controlplane.yaml
➜  code mv talosconfig talos-test/talosconfig            
➜  code mv worker.yaml talos-test/worker.yaml
➜  code ls
```

```
talosctl validate --config talos-test/controlplane.yaml --mode cloud
talosctl validate --config talos-test/worker.yaml --mode cloud
```

```
CP_COUNT=1
while [[ "$CP_COUNT" -lt 4 ]]; do
  aws --profile development ec2 run-instances \
    --region $REGION \
    --image-id $AMI \
    --count 1 \
    --instance-type t3.small \
    --user-data file://controlplane.yaml \
    --subnet-id $SUBNET \
    --security-group-ids $SECURITY_GROUP \
    --associate-public-ip-address \
    --tag-specifications "ResourceType=instance,Tags=[{Key=Name,Value=talos-aws-tutorial-cp-$CP_COUNT}]"
  ((CP_COUNT++))
done
```

```
  code CP_COUNT=1
while [[ "$CP_COUNT" -lt 4 ]]; do
  aws --profile development ec2 run-instances \
    --region $REGION \
    --image-id $AMI \
    --count 1 \
    --instance-type t3.small \
    --user-data file://controlplane.yaml \
    --subnet-id $SUBNET \
    --security-group-ids $SECURITY_GROUP \
    --associate-public-ip-address \
    --tag-specifications "ResourceType=instance,Tags=[{Key=Name,Value=talos-aws-tutorial-cp-$CP_COUNT}]"
  ((CP_COUNT++))
done
{
    "Groups": [],
    "Instances": [
        {
            "AmiLaunchIndex": 0,
            "ImageId": "ami-09360283b6eec5d54",
            "InstanceId": "i-02b3e6232c3039676",
            "InstanceType": "t3.small",
            "LaunchTime": "2024-03-05T23:07:27+00:00",
            "Monitoring": {
                "State": "disabled"
            },
            "Placement": {
                "AvailabilityZone": "us-east-1a",
                "GroupName": "",
                "Tenancy": "default"
            },
            "PrivateDnsName": "ip-10-200-1-90.ec2.internal",
            "PrivateIpAddress": "10.200.1.90",
            "ProductCodes": [],
            "PublicDnsName": "",
            "State": {
                "Code": 0,
                "Name": "pending"
            },
            "StateTransitionReason": "",
            "SubnetId": "subnet-0bdf61de47f04f337",
            "VpcId": "vpc-0d2df8f32c707fd36",
            "Architecture": "x86_64",
            "BlockDeviceMappings": [],
            "ClientToken": "e07479b0-456b-4d82-9ae5-6580ced62765",
            "EbsOptimized": false,
            "EnaSupport": true,
            "Hypervisor": "xen",
            "NetworkInterfaces": [
                {
                    "Attachment": {
                        "AttachTime": "2024-03-05T23:07:27+00:00",
                        "AttachmentId": "eni-attach-04766c138cbe5a4d5",
                        "DeleteOnTermination": true,
                        "DeviceIndex": 0,
                        "Status": "attaching",
                        "NetworkCardIndex": 0
                    },
                    "Description": "",
                    "Groups": [
                        {
                            "GroupName": "talos-aws-tutorial-sg",
                            "GroupId": "sg-01e3966613a73846c"
                        }
                    ],
                    "Ipv6Addresses": [],
                    "MacAddress": "02:4f:cb:de:5d:1d",
                    "NetworkInterfaceId": "eni-0934a6e4fb38d2585",
                    "OwnerId": "339735964233",
                    "PrivateDnsName": "ip-10-200-1-90.ec2.internal",
                    "PrivateIpAddress": "10.200.1.90",
                    "PrivateIpAddresses": [
                        {
                            "Primary": true,
                            "PrivateDnsName": "ip-10-200-1-90.ec2.internal",
                            "PrivateIpAddress": "10.200.1.90"
                        }
                    ],
                    "SourceDestCheck": true,
                    "Status": "in-use",
                    "SubnetId": "subnet-0bdf61de47f04f337",
                    "VpcId": "vpc-0d2df8f32c707fd36",
                    "InterfaceType": "interface"
                }
            ],
            "RootDeviceName": "/dev/xvda",
            "RootDeviceType": "ebs",
            "SecurityGroups": [
                {
                    "GroupName": "talos-aws-tutorial-sg",
                    "GroupId": "sg-01e3966613a73846c"
                }
            ],
            "SourceDestCheck": true,
            "StateReason": {
                "Code": "pending",
                "Message": "pending"
            },
            "Tags": [
                {
                    "Key": "Name",
                    "Value": "talos-aws-tutorial-cp-1"
                }
            ],
            "VirtualizationType": "hvm",
            "CpuOptions": {
                "CoreCount": 1,
                "ThreadsPerCore": 2
            },
            "CapacityReservationSpecification": {
                "CapacityReservationPreference": "open"
            },
            "MetadataOptions": {
                "State": "pending",
                "HttpTokens": "required",
                "HttpPutResponseHopLimit": 2,
                "HttpEndpoint": "enabled",
                "HttpProtocolIpv6": "disabled",
                "InstanceMetadataTags": "disabled"
            },
            "EnclaveOptions": {
                "Enabled": false
            },
            "PrivateDnsNameOptions": {
                "HostnameType": "ip-name",
                "EnableResourceNameDnsARecord": false,
                "EnableResourceNameDnsAAAARecord": false
            },
            "MaintenanceOptions": {
                "AutoRecovery": "default"
            },
            "CurrentInstanceBootMode": "legacy-bios"
        }
    ],
    "OwnerId": "339735964233",
    "ReservationId": "r-0c33f672082c7c4b7"
}
{
    "Groups": [],
    "Instances": [
        {
            "AmiLaunchIndex": 0,
            "ImageId": "ami-09360283b6eec5d54",
            "InstanceId": "i-05f78bb88cad40162",
            "InstanceType": "t3.small",
            "LaunchTime": "2024-03-05T23:07:29+00:00",
            "Monitoring": {
                "State": "disabled"
            },
            "Placement": {
                "AvailabilityZone": "us-east-1a",
                "GroupName": "",
                "Tenancy": "default"
            },
            "PrivateDnsName": "ip-10-200-1-28.ec2.internal",
            "PrivateIpAddress": "10.200.1.28",
            "ProductCodes": [],
            "PublicDnsName": "",
            "State": {
                "Code": 0,
                "Name": "pending"
            },
            "StateTransitionReason": "",
            "SubnetId": "subnet-0bdf61de47f04f337",
            "VpcId": "vpc-0d2df8f32c707fd36",
            "Architecture": "x86_64",
            "BlockDeviceMappings": [],
            "ClientToken": "f0269ce3-b2f7-48d5-89a1-b67a11c256b9",
            "EbsOptimized": false,
            "EnaSupport": true,
            "Hypervisor": "xen",
            "NetworkInterfaces": [
                {
                    "Attachment": {
                        "AttachTime": "2024-03-05T23:07:29+00:00",
                        "AttachmentId": "eni-attach-0da8d16dc2425c17e",
                        "DeleteOnTermination": true,
                        "DeviceIndex": 0,
                        "Status": "attaching",
                        "NetworkCardIndex": 0
                    },
                    "Description": "",
                    "Groups": [
                        {
                            "GroupName": "talos-aws-tutorial-sg",
                            "GroupId": "sg-01e3966613a73846c"
                        }
                    ],
                    "Ipv6Addresses": [],
                    "MacAddress": "02:22:1f:53:7b:7f",
                    "NetworkInterfaceId": "eni-0fd2fefab8195026e",
                    "OwnerId": "339735964233",
                    "PrivateDnsName": "ip-10-200-1-28.ec2.internal",
                    "PrivateIpAddress": "10.200.1.28",
                    "PrivateIpAddresses": [
                        {
                            "Primary": true,
                            "PrivateDnsName": "ip-10-200-1-28.ec2.internal",
                            "PrivateIpAddress": "10.200.1.28"
                        }
                    ],
                    "SourceDestCheck": true,
                    "Status": "in-use",
                    "SubnetId": "subnet-0bdf61de47f04f337",
                    "VpcId": "vpc-0d2df8f32c707fd36",
                    "InterfaceType": "interface"
                }
            ],
            "RootDeviceName": "/dev/xvda",
            "RootDeviceType": "ebs",
            "SecurityGroups": [
                {
                    "GroupName": "talos-aws-tutorial-sg",
                    "GroupId": "sg-01e3966613a73846c"
                }
            ],
            "SourceDestCheck": true,
            "StateReason": {
                "Code": "pending",
                "Message": "pending"
            },
            "Tags": [
                {
                    "Key": "Name",
                    "Value": "talos-aws-tutorial-cp-2"
                }
            ],
            "VirtualizationType": "hvm",
            "CpuOptions": {
                "CoreCount": 1,
                "ThreadsPerCore": 2
            },
            "CapacityReservationSpecification": {
                "CapacityReservationPreference": "open"
            },
            "MetadataOptions": {
                "State": "pending",
                "HttpTokens": "required",
                "HttpPutResponseHopLimit": 2,
                "HttpEndpoint": "enabled",
                "HttpProtocolIpv6": "disabled",
                "InstanceMetadataTags": "disabled"
            },
            "EnclaveOptions": {
                "Enabled": false
            },
            "PrivateDnsNameOptions": {
                "HostnameType": "ip-name",
                "EnableResourceNameDnsARecord": false,
                "EnableResourceNameDnsAAAARecord": false
            },
            "MaintenanceOptions": {
                "AutoRecovery": "default"
            },
            "CurrentInstanceBootMode": "legacy-bios"
        }
    ],
    "OwnerId": "339735964233",
    "ReservationId": "r-01fddc9dbccb70a5a"
}
{
    "Groups": [],
    "Instances": [
        {
            "AmiLaunchIndex": 0,
            "ImageId": "ami-09360283b6eec5d54",
            "InstanceId": "i-00cac9069a97504cd",
            "InstanceType": "t3.small",
            "LaunchTime": "2024-03-05T23:07:31+00:00",
            "Monitoring": {
                "State": "disabled"
            },
            "Placement": {
                "AvailabilityZone": "us-east-1a",
                "GroupName": "",
                "Tenancy": "default"
            },
            "PrivateDnsName": "ip-10-200-1-217.ec2.internal",
            "PrivateIpAddress": "10.200.1.217",
            "ProductCodes": [],
            "PublicDnsName": "",
            "State": {
                "Code": 0,
                "Name": "pending"
            },
            "StateTransitionReason": "",
            "SubnetId": "subnet-0bdf61de47f04f337",
            "VpcId": "vpc-0d2df8f32c707fd36",
            "Architecture": "x86_64",
            "BlockDeviceMappings": [],
            "ClientToken": "ef0c40f6-1072-4752-be1a-2477737d8f75",
            "EbsOptimized": false,
            "EnaSupport": true,
            "Hypervisor": "xen",
            "NetworkInterfaces": [
                {
                    "Attachment": {
                        "AttachTime": "2024-03-05T23:07:31+00:00",
                        "AttachmentId": "eni-attach-0a16e88388e41a541",
                        "DeleteOnTermination": true,
                        "DeviceIndex": 0,
                        "Status": "attaching",
                        "NetworkCardIndex": 0
                    },
                    "Description": "",
                    "Groups": [
                        {
                            "GroupName": "talos-aws-tutorial-sg",
                            "GroupId": "sg-01e3966613a73846c"
                        }
                    ],
                    "Ipv6Addresses": [],
                    "MacAddress": "02:79:86:92:c9:cd",
                    "NetworkInterfaceId": "eni-09a7829c44c1de6a0",
                    "OwnerId": "339735964233",
                    "PrivateDnsName": "ip-10-200-1-217.ec2.internal",
                    "PrivateIpAddress": "10.200.1.217",
                    "PrivateIpAddresses": [
                        {
                            "Primary": true,
                            "PrivateDnsName": "ip-10-200-1-217.ec2.internal",
                            "PrivateIpAddress": "10.200.1.217"
                        }
                    ],
                    "SourceDestCheck": true,
                    "Status": "in-use",
                    "SubnetId": "subnet-0bdf61de47f04f337",
                    "VpcId": "vpc-0d2df8f32c707fd36",
                    "InterfaceType": "interface"
                }
            ],
            "RootDeviceName": "/dev/xvda",
            "RootDeviceType": "ebs",
            "SecurityGroups": [
                {
                    "GroupName": "talos-aws-tutorial-sg",
                    "GroupId": "sg-01e3966613a73846c"
                }
            ],
            "SourceDestCheck": true,
            "StateReason": {
                "Code": "pending",
                "Message": "pending"
            },
            "Tags": [
                {
                    "Key": "Name",
                    "Value": "talos-aws-tutorial-cp-3"
                }
            ],
            "VirtualizationType": "hvm",
            "CpuOptions": {
                "CoreCount": 1,
                "ThreadsPerCore": 2
            },
            "CapacityReservationSpecification": {
                "CapacityReservationPreference": "open"
            },
            "MetadataOptions": {
                "State": "pending",
                "HttpTokens": "required",
                "HttpPutResponseHopLimit": 2,
                "HttpEndpoint": "enabled",
                "HttpProtocolIpv6": "disabled",
                "InstanceMetadataTags": "disabled"
            },
            "EnclaveOptions": {
                "Enabled": false
            },
            "PrivateDnsNameOptions": {
                "HostnameType": "ip-name",
                "EnableResourceNameDnsARecord": false,
                "EnableResourceNameDnsAAAARecord": false
            },
            "MaintenanceOptions": {
                "AutoRecovery": "default"
            },
            "CurrentInstanceBootMode": "legacy-bios"
        }
    ],
    "OwnerId": "339735964233",
    "ReservationId": "r-04ddf4dfbf96784f9"
}
➜  code 






````


The next step of the onboarding guide is to create the worker nodes but I'm going to stop. And see if I can set that up in a different way using asg+event+lambda situation. I'll finish the rest of the tutorial first but I shouldn't need the worker nodes yet. 

```
CP_NODE_1_IP=10.200.1.90
CP_NODE_2_IP=10.200.1.28
CP_NODE_3_IP=10.200.1.217
```

```
aws --profile development elbv2 register-targets \
    --region $REGION \
    --target-group-arn $TARGET_GROUP_ARN \
    --targets Id=$CP_NODE_1_IP  Id=$CP_NODE_2_IP  Id=$CP_NODE_3_IP
```

```
➜  code aws --profile development elbv2 create-listener \
    --region $REGION \
    --load-balancer-arn $LOAD_BALANCER_ARN \
    --protocol TCP \
    --port 443 \
    --default-actions Type=forward,TargetGroupArn=$TARGET_GROUP_ARN
{
    "Listeners": [
        {
            "ListenerArn": "arn:aws:elasticloadbalancing:us-east-1:339735964233:listener/net/talos-aws-tutorial-lb/eedbb2eaa739616a/c9cd37340bc11237",
            "LoadBalancerArn": "arn:aws:elasticloadbalancing:us-east-1:339735964233:loadbalancer/net/talos-aws-tutorial-lb/eedbb2eaa739616a",
            "Port": 443,
            "Protocol": "TCP",
            "DefaultActions": [
                {
                    "Type": "forward",
                    "TargetGroupArn": "arn:aws:elasticloadbalancing:us-east-1:339735964233:targetgroup/talos-aws-tutorial-tg/7146ea36aa0472cc",
                    "ForwardConfig": {
                        "TargetGroups": [
                            {
                                "TargetGroupArn": "arn:aws:elasticloadbalancing:us-east-1:339735964233:targetgroup/talos-aws-tutorial-tg/7146ea36aa0472cc"
                            }
                        ]
                    }
                }
            ]
        }
    ]
}
```

I grabbed one of the control plane instance public ips ipv4: 3.234.177.183

```
CP_PUBLIC_IP_1=3.234.177.183
```

```
➜  code CP_PUBLIC_IP_1=3.234.177.183
➜  code talosctl --talosconfig talosconfig config endpoint $CP_PUBLIC_IP_1
➜  code talosctl --talosconfig talosconfig config node $CP_PUBLIC_IP_1
➜  code talosctl --talosconfig talosconfig bootstrap
```

```
talosctl --talosconfig talosconfig kubeconfig .
```

Awesome it worked...
```
➜  code talosctl --talosconfig talosconfig  health

discovered nodes: ["10.200.1.217" "10.200.1.28" "10.200.1.90"]
waiting for etcd to be healthy: ...
waiting for etcd to be healthy: OK
waiting for etcd members to be consistent across nodes: ...
waiting for etcd members to be consistent across nodes: OK
waiting for etcd members to be control plane nodes: ...
waiting for etcd members to be control plane nodes: OK
waiting for apid to be ready: ...
waiting for apid to be ready: OK
waiting for all nodes memory sizes: ...
waiting for all nodes memory sizes: OK
waiting for all nodes disk sizes: ...
waiting for all nodes disk sizes: OK
waiting for kubelet to be healthy: ...
waiting for kubelet to be healthy: OK
waiting for all nodes to finish boot sequence: ...
waiting for all nodes to finish boot sequence: OK
waiting for all k8s nodes to report: ...
waiting for all k8s nodes to report: can't find expected node with IPs ["10.200.1.217"]
waiting for all k8s nodes to report: OK
waiting for all k8s nodes to report ready: ...
waiting for all k8s nodes to report ready: some nodes are not ready: [ip-10-200-1-217 ip-10-200-1-28 ip-10-200-1-90]
waiting for all k8s nodes to report ready: some nodes are not ready: [ip-10-200-1-90]
waiting for all k8s nodes to report ready: OK
waiting for all control plane static pods to be running: ...
waiting for all control plane static pods to be running: OK
waiting for all control plane components to be ready: ...
waiting for all control plane components to be ready: can't find expected node with IPs ["10.200.1.28"]
waiting for all control plane components to be ready: expected number of pods for kube-apiserver to be 3, got 2
waiting for all control plane components to be ready: OK
waiting for kube-proxy to report ready: ...
waiting for kube-proxy to report ready: OK
waiting for coredns to report ready: ...
waiting for coredns to report ready: OK
waiting for all k8s nodes to report schedulable: ...
waiting for all k8s nodes to report schedulable: OK
```