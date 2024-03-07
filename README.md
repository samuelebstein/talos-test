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

Using the kubeconfig that was created...

```
➜  talos-test git:(main) ✗ kubectl --kubeconfig=kubeconfig get nodes -o wide
NAME              STATUS   ROLES           AGE     VERSION   INTERNAL-IP    EXTERNAL-IP   OS-IMAGE         KERNEL-VERSION   CONTAINER-RUNTIME
ip-10-200-1-217   Ready    control-plane   5m30s   v1.29.2   10.200.1.217   <none>        Talos (v1.6.4)   6.1.74-talos     containerd://1.7.13
ip-10-200-1-28    Ready    control-plane   5m31s   v1.29.2   10.200.1.28    <none>        Talos (v1.6.4)   6.1.74-talos     containerd://1.7.13
ip-10-200-1-90    Ready    control-plane   5m30s   v1.29.2   10.200.1.90    <none>        Talos (v1.6.4)   6.1.74-talos     containerd://1.7.13
```

```
➜  talos-test git:(main) ✗ talosctl --talosconfig talosconfig service       
NODE            SERVICE      STATE     HEALTH   LAST CHANGE   LAST EVENT
3.234.177.183   apid         Running   OK       31m13s ago    Health check successful
3.234.177.183   containerd   Running   OK       31m20s ago    Health check successful
3.234.177.183   cri          Running   OK       31m15s ago    Health check successful
3.234.177.183   dashboard    Running   ?        31m19s ago    Process Process(["/sbin/dashboard"]) started with PID 1305
3.234.177.183   etcd         Running   OK       13m24s ago    Health check successful
3.234.177.183   kubelet      Running   OK       31m3s ago     Health check successful
3.234.177.183   machined     Running   OK       31m26s ago    Health check successful
3.234.177.183   trustd       Running   OK       31m15s ago    Health check successful
3.234.177.183   udevd        Running   OK       31m24s ago    Health check successful
```

Now lets try to bring up some worker nodes..

Create the asg, launch config, try to bring up a worker without user data. it should succeed and then we can terminate it.

```
aws --profile development autoscaling create-launch-configuration \
    --launch-configuration-name talos-aws-tutorial-worker-launch-config \
    --image-id $AMI \
    --instance-type t3.small \
    --security-groups $SECURITY_GROUP # Using same security group as the control plane nodes. no user data
    # --iam-instance-profile your-iam-role # Maybe add this later 
```
```
aws --profile development ec2 create-launch-template \
    --launch-template-name "talos-aws-tutorial-worker-launch-config" \
    --version-description "version1" \
    --launch-template-data "{\"ImageId\":\"$AMI\",\"InstanceType\":\"t3.small\",\"SecurityGroupIds\":[\"$SECURITY_GROUP\"],\"TagSpecifications\":[{\"ResourceType\":\"instance\",\"Tags\":[{\"Key\":\"Purpose\",\"Value\":\"talos-aws-tutorial-worker\"}]}]}"    

```        

```
➜  talos-test git:(main) ✗ aws --profile development ec2 create-launch-template \
    --launch-template-name "talos-aws-tutorial-worker-launch-config" \
    --version-description "version1" \
    --launch-template-data "{\"ImageId\":\"$AMI\",\"InstanceType\":\"t3.small\",\"SecurityGroupIds\":[\"$SECURITY_GROUP\"],\"TagSpecifications\":[{\"ResourceType\":\"instance\",\"Tags\":[{\"Key\":\"Purpose\",\"Value\":\"talos-aws-tutorial-worker\"}]}]}"    
{
    "LaunchTemplate": {
        "LaunchTemplateId": "lt-07671e4a96fe36737",
        "LaunchTemplateName": "talos-aws-tutorial-worker-launch-config",
        "CreateTime": "2024-03-06T00:25:45+00:00",
        "CreatedBy": "arn:aws:sts::339735964233:assumed-role/admin/sam_ebstein",
        "DefaultVersionNumber": 1,
        "LatestVersionNumber": 1
    }
}
```

```
aws --profile development autoscaling create-auto-scaling-group \
    --auto-scaling-group-name talos-workers-asg \
    --launch-configuration-name talos-aws-tutorial-worker-launch-config \
    --min-size 0 \
    --max-size 3 \
    --desired-capacity 1 \
    --vpc-zone-identifier $SUBNET \
    --tags "Key=Name,Value=talos-worker,PropagateAtLaunch=true"
```    

that successfully created an autoscaling group...


Now on to trying to create a lambda function...

```
go mod init github.com/samuelebstein/talos-test/talos-applier-lambda-function

go get github.com/aws/aws-lambda-go/lambda

GOOS=linux GOARCH=amd64 go build -o main

zip function.zip main

```


```
aws --profile development iam create-role --role-name TalosLambdaExecutionRole --assume-role-policy-document file://trust-policy.json

➜  talos-applier-lambda-function git:(main) ✗ aws --profile development iam create-role --role-name TalosLambdaExecutionRole --assume-role-policy-document file://trust-policy.json
{
    "Role": {
        "Path": "/",
        "RoleName": "TalosLambdaExecutionRole",
        "RoleId": "AROAU6GOXNJETQEUZCNPZ",
        "Arn": "arn:aws:iam::339735964233:role/TalosLambdaExecutionRole",
        "CreateDate": "2024-03-06T18:28:59+00:00",
        "AssumeRolePolicyDocument": {
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Effect": "Allow",
                    "Principal": {
                        "Service": "lambda.amazonaws.com"
                    },
                    "Action": "sts:AssumeRole"
                }
            ]
        }
    }
}

```

Lets create the test secret in secrets manager before creating the policies that we'll attaach to the lambda execution role...


```
➜  talos-test git:(main) ✗ aws --profile development secretsmanager create-secret --name "sam-ebstein-test-talosconfig" \
    --description "Talos configuration for sam-ebstein-test" \
    --secret-string file://talosconfig
{
    "ARN": "arn:aws:secretsmanager:us-east-1:339735964233:secret:sam-ebstein-test-talosconfig-CnZCDQ",
    "Name": "sam-ebstein-test-talosconfig",
    "VersionId": "1b5ce0ea-6ebb-4d9d-92a0-c2aa7da4cf19"
}
```

Is the talosconfig something that changes over time?? Because if so, then the secret data would have to updated every time the endpoing (for example) changes..for instance if the ec2 instance falls down


lets create the worker secret also

```
➜  talos-test git:(main) ✗ aws --profile development secretsmanager create-secret --name "sam-ebstein-test-talos-worker-yaml" \
    --description "Talos worker yaml configuration for sam-ebstein-test" \
    --secret-string file://worker.yaml
{
    "ARN": "arn:aws:secretsmanager:us-east-1:339735964233:secret:sam-ebstein-test-talos-worker-yaml-cyRmM4",
    "Name": "sam-ebstein-test-talos-worker-yaml",
    "VersionId": "afd500ae-15cf-480a-b568-3b1bb76b8071"
}
```

Okay now lets go back and create the policies

```
➜  talos-applier-lambda-function git:(main) ✗ aws --profile development iam create-policy --policy-name LambdaExecutionPermissions --policy-document file://permissions-policy.json
{
    "Policy": {
        "PolicyName": "LambdaExecutionPermissions",
        "PolicyId": "ANPAU6GOXNJE3BOEPONB4",
        "Arn": "arn:aws:iam::339735964233:policy/LambdaExecutionPermissions",
        "Path": "/",
        "DefaultVersionId": "v1",
        "AttachmentCount": 0,
        "PermissionsBoundaryUsageCount": 0,
        "IsAttachable": true,
        "CreateDate": "2024-03-06T18:51:55+00:00",
        "UpdateDate": "2024-03-06T18:51:55+00:00"
    }
}
```


```
aws --profile development iam attach-role-policy --role-name TalosLambdaExecutionRole --policy-arn "arn:aws:iam::339735964233:policy/LambdaExecutionPermissions"

```


trying secrets retrieval in first deployment of code...
```

GOOS=linux GOARCH=amd64 go build -o main
zip deployment.zip main


```

```
  talos-applier-lambda-function git:(main) ✗ aws --profile development lambda create-function \                                                                                                                      
    --function-name SamEbsteinTalosLambdaTest \                                    
    --runtime provided.al2 \                                                                                                                                          
    --role arn:aws:iam::339735964233:role/TalosLambdaExecutionRole \                                       
    --handler main \                                                               
    --zip-file fileb://deployment.zip \
    --architecture x86_64 \
    --timeout 15 \                                                                 
    --memory-size 128                                                                                                                                                                      
{                             
    "FunctionName": "SamEbsteinTalosLambdaTest",                      
    "FunctionArn": "arn:aws:lambda:us-east-1:339735964233:function:SamEbsteinTalosLambdaTest",                                                                                                                        
    "Runtime": "provided.al2",
    "Role": "arn:aws:iam::339735964233:role/TalosLambdaExecutionRole",             
    "Handler": "main",       
    "CodeSize": 7233800,                                                           
    "Description": "",                                                             
    "Timeout": 15,                                                                           
    "MemorySize": 128,                                                             
    "LastModified": "2024-03-06T19:12:16.950+0000",                                                        
    "CodeSha256": "Z0SGkZomq2TGtQiLvJ3/DV/Lfyf9fgBUFaxB7YNgsPw=",                                          
    "Version": "$LATEST",
    "TracingConfig": {                                                                       
        "Mode": "PassThrough"
    },                                                                                       
    "RevisionId": "29412239-452c-42f6-a114-e081320c3e53",                                                  
    "State": "Pending",            
    "StateReason": "The function is being created.",                                                       
    "StateReasonCode": "Creating",
    "PackageType": "Zip",                                                                                                                                             
    "Architectures": [   
        "x86_64"             
    ],                                                                                                                                                                
    "EphemeralStorage": {                                                          
        "Size": 512       
    },                             
    "SnapStart": {                                                                 
        "ApplyOn": "None",   
        "OptimizationStatus": "Off"                                                                                                                                                        
    },                                                                             
    "RuntimeVersionConfig": {                 
        "RuntimeVersionArn": "arn:aws:lambda:us-east-1::runtime:e44362e335db9c887e4819f03950e642c889a449eb010a6f1b4cb1a0d7e5c92b"                                                                                     
    },                                                                                       
    "LoggingConfig": {                        
        "LogFormat": "Text",                  
        "LogGroup": "/aws/lambda/SamEbsteinTalosLambdaTest"                                  
    }                                                
}  
```

function is created.. trying to create a test event with this input

{
  "name": "Test User"
}

received this error:

```
{
  "errorType": "Runtime.InvalidEntrypoint",
  "errorMessage": "RequestId: 407ca2c1-4e10-48f1-8594-1ba30a0aab85 Error: Couldn't find valid bootstrap(s): [/var/task/bootstrap /opt/bootstrap]"
}
```

Think its because I used custom runtime on creation

updating the function runtime:

```
➜  talos-applier-lambda-function git:(main) ✗ aws --profile development lambda update-function-configuration \
    --function-name SamEbsteinTalosLambdaTest \
    --runtime go1.x \
    --handler main
{
    "FunctionName": "SamEbsteinTalosLambdaTest",
    "FunctionArn": "arn:aws:lambda:us-east-1:339735964233:function:SamEbsteinTalosLambdaTest",
    "Runtime": "go1.x",
    "Role": "arn:aws:iam::339735964233:role/TalosLambdaExecutionRole",
    "Handler": "main",
    "CodeSize": 7233800,
    "Description": "",
    "Timeout": 15,
    "MemorySize": 128,
    "LastModified": "2024-03-06T19:19:00.000+0000",
    "CodeSha256": "Z0SGkZomq2TGtQiLvJ3/DV/Lfyf9fgBUFaxB7YNgsPw=",
    "Version": "$LATEST",
    "TracingConfig": {
        "Mode": "PassThrough"
    },
    "RevisionId": "c8ab45ae-f461-444f-bba0-fb29a1a1c370",
    "State": "Active",
    "LastUpdateStatus": "InProgress",
    "LastUpdateStatusReason": "The function is being created.",
    "LastUpdateStatusReasonCode": "Creating",
    "PackageType": "Zip",
    "Architectures": [
        "x86_64"
    ],
    "EphemeralStorage": {
        "Size": 512
    },
    "SnapStart": {
        "ApplyOn": "None",
        "OptimizationStatus": "Off"
    },
    "RuntimeVersionConfig": {
        "RuntimeVersionArn": "arn:aws:lambda:us-east-1::runtime:30052276b0b7733e82eddf1f0942de1022c7dfbc0ca93cfc121c868194868dec"
    },
    "LoggingConfig": {
        "LogFormat": "Text",
        "LogGroup": "/aws/lambda/SamEbsteinTalosLambdaTest"
    }
}
```

Function succeeds now.

```
➜  talos-applier-lambda-function git:(main) ✗ aws --profile development events put-rule \
    --name "TalosWorkersAsgScaleUp" \
    --event-pattern file://event-pattern.json \
    --state ENABLED
{
    "RuleArn": "arn:aws:events:us-east-1:339735964233:rule/TalosWorkersAsgScaleUp"
}
```


```
➜  talos-applier-lambda-function git:(main) ✗ aws --profile development events put-targets \
    --rule "TalosWorkersAsgScaleUp" \
    --targets "Id"="1","Arn"="arn:aws:lambda:us-east-1:339735964233:function:SamEbsteinTalosLambdaTest"
{
    "FailedEntryCount": 0,
    "FailedEntries": []
}
```    


Add resource based policy to allow lambda to be invoked:

```
➜  talos-applier-lambda-function git:(main) ✗ aws --profile development lambda add-permission \
    --function-name SamEbsteinTalosLambdaTest \
    --statement-id EventBridgeInvokeSamEbsteinTalosLambdaTest \
    --action lambda:InvokeFunction \
    --principal events.amazonaws.com \
    --source-arn arn:aws:events:us-east-1:339735964233:rule/TalosWorkersAsgScaleUp
{
    "Statement": "{\"Sid\":\"EventBridgeInvokeSamEbsteinTalosLambdaTest\",\"Effect\":\"Allow\",\"Principal\":{\"Service\":\"events.amazonaws.com\"},\"Action\":\"lambda:InvokeFunction\",\"Resource\":\"arn:aws:lambda:us-east-1:339735964233:function:SamEbsteinTalosLambdaTest\",\"Condition\":{\"ArnLike\":{\"AWS:SourceArn\":\"arn:aws:events:us-east-1:339735964233:rule/TalosWorkersAsgScaleUp\"}}}"
}
```    




should be able to bring in this client in the code:

https://github.com/siderolabs/talos/blob/8c79539914324eee64dbdaf1f535fc4e20da55e8/pkg/machinery/client/client.go#L243


updating the policy to include ec2 describe instances

```
➜  talos-applier-lambda-function git:(main) ✗ aws --profile development iam create-policy-version \
    --policy-arn arn:aws:iam::339735964233:policy/LambdaExecutionPermissions \
    --policy-document file://permissions-policy.json \
    --set-as-default
{
    "PolicyVersion": {
        "VersionId": "v2",
        "IsDefaultVersion": true,
        "CreateDate": "2024-03-07T00:45:41+00:00"
    }
}
```

```
➜  talos-applier-lambda-function git:(main) ✗ GOOS=linux GOARCH=amd64 go build -o main

➜  talos-applier-lambda-function git:(main) ✗ zip deployment.zip main                 

  adding: main (deflated 48%)
➜  talos-applier-lambda-function git:(main) ✗ 
```


```
➜  talos-applier-lambda-function git:(main) ✗ aws --profile development lambda update-function-code --function-name SamEbsteinTalosLambdaTest --zip-file fileb://deployment.zip
{
    "FunctionName": "SamEbsteinTalosLambdaTest",
    "FunctionArn": "arn:aws:lambda:us-east-1:339735964233:function:SamEbsteinTalosLambdaTest",
    "Runtime": "go1.x",
    "Role": "arn:aws:iam::339735964233:role/TalosLambdaExecutionRole",
    "Handler": "main",
    "CodeSize": 7370720,
    "Description": "",
    "Timeout": 15,
    "MemorySize": 128,
    "LastModified": "2024-03-07T00:51:40.000+0000",
    "CodeSha256": "p4nIa48M8YisuAL0iCKMu+Nsz73A+eOQZX1cpTfG3hQ=",
    "Version": "$LATEST",
    "TracingConfig": {
        "Mode": "PassThrough"
    },
    "RevisionId": "4ade823f-e34d-4161-b52f-ef17188bbaca",
    "State": "Active",
    "LastUpdateStatus": "InProgress",
    "LastUpdateStatusReason": "The function is being created.",
    "LastUpdateStatusReasonCode": "Creating",
    "PackageType": "Zip",
    "Architectures": [
        "x86_64"
    ],
    "EphemeralStorage": {
        "Size": 512
    },
    "SnapStart": {
        "ApplyOn": "None",
        "OptimizationStatus": "Off"
    },
    "RuntimeVersionConfig": {
        "RuntimeVersionArn": "arn:aws:lambda:us-east-1::runtime:30052276b0b7733e82eddf1f0942de1022c7dfbc0ca93cfc121c868194868dec"
    },
    "LoggingConfig": {
        "LogFormat": "Text",
        "LogGroup": "/aws/lambda/SamEbsteinTalosLambdaTest"
    }
}

```

```
aws --profile development lambda get-function-configuration --function-name SamEbsteinTalosLambdaTest

```

running into a dereferencing error, redeploying the code...

didn't have public ip addresss associated with my autoscaling group..

```
➜  talos-applier-lambda-function git:(main) ✗ aws --profile development ec2 create-launch-template-version \
    --launch-template-name talos-aws-tutorial-worker-launch-config \
    --version-description "version with public IP" \
    --launch-template-data '{"NetworkInterfaces":[{"DeviceIndex":0,"AssociatePublicIpAddress":true}]}'
{
    "LaunchTemplateVersion": {
        "LaunchTemplateId": "lt-07671e4a96fe36737",
        "LaunchTemplateName": "talos-aws-tutorial-worker-launch-config",
        "VersionNumber": 2,
        "VersionDescription": "version with public IP",
        "CreateTime": "2024-03-07T01:20:19+00:00",
        "CreatedBy": "arn:aws:sts::339735964233:assumed-role/admin/sam_ebstein",
        "DefaultVersion": false,
        "LaunchTemplateData": {
            "NetworkInterfaces": [
                {
                    "AssociatePublicIpAddress": true,
                    "DeviceIndex": 0
                }
            ]
        }
    }
}
```    


```
➜  talos-applier-lambda-function git:(main) ✗ aws --profile development ec2 modify-launch-template \
    --launch-template-name talos-aws-tutorial-worker-launch-config \
    --default-version 2
{
    "LaunchTemplate": {
        "LaunchTemplateId": "lt-07671e4a96fe36737",
        "LaunchTemplateName": "talos-aws-tutorial-worker-launch-config",
        "CreateTime": "1970-01-01T00:00:00+00:00",
        "CreatedBy": "arn:aws:sts::339735964233:assumed-role/admin/sam_ebstein",
        "DefaultVersionNumber": 2,
        "LatestVersionNumber": 2
    }
}
```

```    
➜  talos-applier-lambda-function git:(main) ✗ aws --profile development ec2 create-launch-template-version \
    --launch-template-name talos-aws-tutorial-worker-launch-config \
    --version-description "version with public IP and other params" \
    --launch-template-data "{\"ImageId\":\"$AMI\",\"InstanceType\":\"t3.small\",\"SecurityGroupIds\":[\"$SECURITY_GROUP\"],\"TagSpecifications\":[{\"ResourceType\":\"instance\",\"Tags\":[{\"Key\":\"Purpose\",\"Value\":\"talos-aws-tutorial-worker\"}]}],\"NetworkInterfaces\":[{\"DeviceIndex\":0,\"AssociatePublicIpAddress\":true}]}"
{
    "LaunchTemplateVersion": {
        "LaunchTemplateId": "lt-07671e4a96fe36737",
        "LaunchTemplateName": "talos-aws-tutorial-worker-launch-config",
        "VersionNumber": 3,
        "VersionDescription": "version with public IP and other params",
        "CreateTime": "2024-03-07T01:26:39+00:00",
        "CreatedBy": "arn:aws:sts::339735964233:assumed-role/admin/sam_ebstein",
        "DefaultVersion": false,
        "LaunchTemplateData": {
            "NetworkInterfaces": [
                {
                    "AssociatePublicIpAddress": true,
                    "DeviceIndex": 0
                }
            ],
            "ImageId": "ami-09360283b6eec5d54",
            "InstanceType": "t3.small",
            "TagSpecifications": [
                {
                    "ResourceType": "instance",
                    "Tags": [
                        {
                            "Key": "Purpose",
                            "Value": "talos-aws-tutorial-worker"
                        }
                    ]
                }
            ],
            "SecurityGroupIds": [
                "sg-01e3966613a73846c"
            ]
        }
    },
    "Warning": {
        "Errors": [
            {
                "Code": "InvalidParameterCombination",
                "Message": "If you specify a network interface, you must specify all security groups as part of the network interface."
            }
        ]
    }
}
```
```
➜  talos-applier-lambda-function git:(main) ✗ aws --profile development ec2 modify-launch-template \
    --launch-template-name talos-aws-tutorial-worker-launch-config \
    --default-version 3
{
    "LaunchTemplate": {
        "LaunchTemplateId": "lt-07671e4a96fe36737",
        "LaunchTemplateName": "talos-aws-tutorial-worker-launch-config",
        "CreateTime": "1970-01-01T00:00:00+00:00",
        "CreatedBy": "arn:aws:sts::339735964233:assumed-role/admin/sam_ebstein",
        "DefaultVersionNumber": 3,
        "LatestVersionNumber": 3
    }
}

```    
    

I deleted the autoscaling group to setup with the launch template instead...

creating new launch template version due to security group issues

```
➜  talos-applier-lambda-function git:(main) ✗ aws --profile development ec2 create-launch-template-version \
    --launch-template-name talos-aws-tutorial-worker-launch-config \
    --version-description "version with public IP and other params" \
    --launch-template-data "{\"ImageId\":\"$AMI\",\"InstanceType\":\"t3.small\",\"TagSpecifications\":[{\"ResourceType\":\"instance\",\"Tags\":[{\"Key\":\"Purpose\",\"Value\":\"talos-aws-tutorial-worker\"}]}],\"NetworkInterfaces\":[{\"DeviceIndex\":0,\"AssociatePublicIpAddress\":true, \"Groups\":[\"$SECURITY_GROUP\"]}]}"
{
    "LaunchTemplateVersion": {
        "LaunchTemplateId": "lt-07671e4a96fe36737",
        "LaunchTemplateName": "talos-aws-tutorial-worker-launch-config",
        "VersionNumber": 4,
        "VersionDescription": "version with public IP and other params",
        "CreateTime": "2024-03-07T01:39:57+00:00",
        "CreatedBy": "arn:aws:sts::339735964233:assumed-role/admin/sam_ebstein",
        "DefaultVersion": false,
        "LaunchTemplateData": {
            "NetworkInterfaces": [
                {
                    "AssociatePublicIpAddress": true,
                    "DeviceIndex": 0,
                    "Groups": [
                        "sg-01e3966613a73846c"
                    ]
                }
            ],
            "ImageId": "ami-09360283b6eec5d54",
            "InstanceType": "t3.small",
            "TagSpecifications": [
                {
                    "ResourceType": "instance",
                    "Tags": [
                        {
                            "Key": "Purpose",
                            "Value": "talos-aws-tutorial-worker"
                        }
                    ]
                }
            ]
        }
    }
}
```

```
➜  talos-applier-lambda-function git:(main) ✗ aws --profile development ec2 modify-launch-template \
    --launch-template-name talos-aws-tutorial-worker-launch-config \
    --default-version 4
{
    "LaunchTemplate": {
        "LaunchTemplateId": "lt-07671e4a96fe36737",
        "LaunchTemplateName": "talos-aws-tutorial-worker-launch-config",
        "CreateTime": "1970-01-01T00:00:00+00:00",
        "CreatedBy": "arn:aws:sts::339735964233:assumed-role/admin/sam_ebstein",
        "DefaultVersionNumber": 4,
        "LatestVersionNumber": 4
    }
}
➜  talos-applier-lambda-function git:(main) ✗ aws --profile development autoscaling create-auto-scaling-group \
    --auto-scaling-group-name talos-workers-asg \
    --launch-template "LaunchTemplateName=talos-aws-tutorial-worker-launch-config,Version=4" \
    --min-size 0 \
    --max-size 3 \
    --desired-capacity 1 \
    --vpc-zone-identifier $SUBNET \
    --tags "Key=Name,Value=talos-worker,PropagateAtLaunch=true"
```    

okay I've created a new launch template and main.go can find the public ip address.. finally..!

next step is to update the function to get the secret and try to apply the config using the talos go client...!

update function:
```
aws --profile development lambda update-function-code --function-name SamEbsteinTalosLambdaTest --zip-file fileb://deployment.zip
```

Recieved this error:

Mode details (should contain info about response):  Applied configuration with a reboot: this config change can't be applied in immediate mode

also mode details logs too much information...

Lets retry creating the cluster...

delete load balancer (arn)
delete target group (arn)
delete instances (cp ip addresses)
delete listener

➜  talos-test git:(main) ✗ aws --profile development elbv2 create-load-balancer \
    --region $REGION \
    --name talos-aws-tutorial-lb \
    --type network --subnets $SUBNET
{
    "LoadBalancers": [
        {
            "LoadBalancerArn": "arn:aws:elasticloadbalancing:us-east-1:339735964233:loadbalancer/net/talos-aws-tutorial-lb/c0674df5bf28aa94",
            "DNSName": "talos-aws-tutorial-lb-c0674df5bf28aa94.elb.us-east-1.amazonaws.com",
            "CanonicalHostedZoneId": "Z26RNL4JYFTOTI",
            "CreatedTime": "2024-03-07T19:18:28.058000+00:00",
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

➜  talos-test git:(main) ✗ aws --profile development elbv2 create-target-group --region $REGION \ 
    --name talos-aws-tutorial-tg \
    --protocol TCP \
    --port 6443 \
    --target-type ip \
    --vpc-id $VPC 
{
    "TargetGroups": [
        {
            "TargetGroupArn": "arn:aws:elasticloadbalancing:us-east-1:339735964233:targetgroup/talos-aws-tutorial-tg/d7d7505606efd8d1",
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
talosctl gen config talos-k8s-aws-tutorial https://talos-aws-tutorial-lb-c0674df5bf28aa94.elb.us-east-1.amazonaws.com:443 --with-examples=false --with-docs=false
```

```
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
            "InstanceId": "i-07d5db6f33dfca74a",
            "InstanceType": "t3.small",
            "LaunchTime": "2024-03-07T19:25:23+00:00",
            "Monitoring": {
                "State": "disabled"
            },
            "Placement": {
                "AvailabilityZone": "us-east-1a",
                "GroupName": "",
                "Tenancy": "default"
            },
            "PrivateDnsName": "ip-10-200-1-211.ec2.internal",
            "PrivateIpAddress": "10.200.1.211",
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
            "ClientToken": "dcf8a489-b30b-4123-93e4-096c33019a08",
            "EbsOptimized": false,
            "EnaSupport": true,
            "Hypervisor": "xen",
            "NetworkInterfaces": [
                {
                    "Attachment": {
                        "AttachTime": "2024-03-07T19:25:23+00:00",
                        "AttachmentId": "eni-attach-0343444f77318012b",
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
                    "MacAddress": "02:b8:4e:ce:f4:67",
                    "NetworkInterfaceId": "eni-0411b8d486deeb786",
                    "OwnerId": "339735964233",
                    "PrivateDnsName": "ip-10-200-1-211.ec2.internal",
                    "PrivateIpAddress": "10.200.1.211",
                    "PrivateIpAddresses": [
                        {
                            "Primary": true,
                            "PrivateDnsName": "ip-10-200-1-211.ec2.internal",
                            "PrivateIpAddress": "10.200.1.211"
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
    "ReservationId": "r-041efaf73512c935f"
}
{
    "Groups": [],
    "Instances": [
        {
            "AmiLaunchIndex": 0,
            "ImageId": "ami-09360283b6eec5d54",
            "InstanceId": "i-0c4ae87b50f877a6b",
            "InstanceType": "t3.small",
            "LaunchTime": "2024-03-07T19:25:25+00:00",
            "Monitoring": {
                "State": "disabled"
            },
            "Placement": {
                "AvailabilityZone": "us-east-1a",
                "GroupName": "",
                "Tenancy": "default"
            },
            "PrivateDnsName": "ip-10-200-1-22.ec2.internal",
            "PrivateIpAddress": "10.200.1.22",
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
            "ClientToken": "d8b5e9af-1925-4b48-9631-8e82461edea6",
            "EbsOptimized": false,
            "EnaSupport": true,
            "Hypervisor": "xen",
            "NetworkInterfaces": [
                {
                    "Attachment": {
                        "AttachTime": "2024-03-07T19:25:25+00:00",
                        "AttachmentId": "eni-attach-0bcdaf49f1cc190a7",
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
                    "MacAddress": "02:6e:fc:94:19:cb",
                    "NetworkInterfaceId": "eni-05e5435e0acc39e0c",
                    "OwnerId": "339735964233",
                    "PrivateDnsName": "ip-10-200-1-22.ec2.internal",
                    "PrivateIpAddress": "10.200.1.22",
                    "PrivateIpAddresses": [
                        {
                            "Primary": true,
                            "PrivateDnsName": "ip-10-200-1-22.ec2.internal",
                            "PrivateIpAddress": "10.200.1.22"
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
    "ReservationId": "r-04072be168f784592"
}
{
    "Groups": [],
    "Instances": [
        {
            "AmiLaunchIndex": 0,
            "ImageId": "ami-09360283b6eec5d54",
            "InstanceId": "i-0f3ea6ba56f6d043c",
            "InstanceType": "t3.small",
            "LaunchTime": "2024-03-07T19:25:27+00:00",
            "Monitoring": {
                "State": "disabled"
            },
            "Placement": {
                "AvailabilityZone": "us-east-1a",
                "GroupName": "",
                "Tenancy": "default"
            },
            "PrivateDnsName": "ip-10-200-1-252.ec2.internal",
            "PrivateIpAddress": "10.200.1.252",
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
            "ClientToken": "1b566092-1593-4ac3-b9c4-287d1390cf9d",
            "EbsOptimized": false,
            "EnaSupport": true,
            "Hypervisor": "xen",
            "NetworkInterfaces": [
                {
                    "Attachment": {
                        "AttachTime": "2024-03-07T19:25:27+00:00",
                        "AttachmentId": "eni-attach-06634faa6fdf32ba0",
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
                    "MacAddress": "02:d3:84:0c:33:ed",
                    "NetworkInterfaceId": "eni-069d9f383e9ba2285",
                    "OwnerId": "339735964233",
                    "PrivateDnsName": "ip-10-200-1-252.ec2.internal",
                    "PrivateIpAddress": "10.200.1.252",
                    "PrivateIpAddresses": [
                        {
                            "Primary": true,
                            "PrivateDnsName": "ip-10-200-1-252.ec2.internal",
                            "PrivateIpAddress": "10.200.1.252"
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
    "ReservationId": "r-0dd070a903aa451bc"
}
```



```

CP_NODE_1_IP=10.200.1.252
CP_NODE_2_IP=10.200.1.22
CP_NODE_3_IP=10.200.1.211

aws --profile development elbv2 register-targets \
    --region $REGION \
    --target-group-arn $TARGET_GROUP_ARN \
    --targets Id=$CP_NODE_1_IP  Id=$CP_NODE_2_IP  Id=$CP_NODE_3_IP
```
```
➜  talos-applier-lambda-function git:(main) ✗ aws --profile development elbv2 create-listener \
    --region $REGION \
    --load-balancer-arn $LOAD_BALANCER_ARN \
    --protocol TCP \
    --port 443 \
    --default-actions Type=forward,TargetGroupArn=$TARGET_GROUP_ARN
{
    "Listeners": [
        {
            "ListenerArn": "arn:aws:elasticloadbalancing:us-east-1:339735964233:listener/net/talos-aws-tutorial-lb/c0674df5bf28aa94/a75a92a21a9cc10e",
            "LoadBalancerArn": "arn:aws:elasticloadbalancing:us-east-1:339735964233:loadbalancer/net/talos-aws-tutorial-lb/c0674df5bf28aa94",
            "Port": 443,
            "Protocol": "TCP",
            "DefaultActions": [
                {
                    "Type": "forward",
                    "TargetGroupArn": "arn:aws:elasticloadbalancing:us-east-1:339735964233:targetgroup/talos-aws-tutorial-tg/d7d7505606efd8d1",
                    "ForwardConfig": {
                        "TargetGroups": [
                            {
                                "TargetGroupArn": "arn:aws:elasticloadbalancing:us-east-1:339735964233:targetgroup/talos-aws-tutorial-tg/d7d7505606efd8d1"
                            }
                        ]
                    }
                }
            ]
        }
    ]
}
```

```
➜  talos-applier-lambda-function git:(main) ✗ talosctl --talosconfig talosconfig config endpoint 3.239.168.158
➜  talos-applier-lambda-function git:(main) ✗ talosctl --talosconfig talosconfig config node 3.239.168.158
➜  talos-applier-lambda-function git:(main) ✗ talosctl --talosconfig talosconfig bootstrap

➜  talos-applier-lambda-function git:(main) ✗ 
```

Its up now:

➜  talos-applier-lambda-function git:(main) ✗ talosctl --talosconfig talosconfig  health

discovered nodes: ["10.200.1.211" "10.200.1.22" "10.200.1.252"]
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
waiting for all k8s nodes to report: can't find expected node with IPs ["10.200.1.211"]
waiting for all k8s nodes to report: OK
waiting for all k8s nodes to report ready: ...
waiting for all k8s nodes to report ready: OK
waiting for all control plane static pods to be running: ...
waiting for all control plane static pods to be running: OK
waiting for all control plane components to be ready: ...
waiting for all control plane components to be ready: can't find expected node with IPs ["10.200.1.211"]
waiting for all control plane components to be ready: expected number of pods for kube-controller-manager to be 3, got 2
waiting for all control plane components to be ready: OK
waiting for kube-proxy to report ready: ...
waiting for kube-proxy to report ready: OK
waiting for coredns to report ready: ...
waiting for coredns to report ready: OK
waiting for all k8s nodes to report schedulable: ...
waiting for all k8s nodes to report schedulable: OK


```
aws --profile development secretsmanager put-secret-value \
    --secret-id "sam-ebstein-test-talosconfig" \
    --secret-string file://talosconfig


aws --profile development secretsmanager put-secret-value \
    --secret-id "sam-ebstein-test-talos-worker-yaml" \
    --secret-string file://worker.yaml    
```

