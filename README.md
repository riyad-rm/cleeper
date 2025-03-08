# Cleeper 🚀

Cleeper is a tool designed to **automatically shut down resources within AWS** to optimize costs and management.\
It currently supports:

- **EC2 Instances** → Stops the instances.
- **RDS Databases** → Stops RDS instances and clusters.
- **ASG Auto Scaling Groups** → Terminates instances within the ASG and suspends the launch process.

---

## 📌 How to Deploy

### **Prerequisites**

Make sure you have the following installed:

- **Terraform**
- **Go**
- **zip**

### **Deployment Steps**

1. Navigate to the `deploy/` folder.
2. Run the `deploy.sh` script:
   ```bash
   ./deploy.sh
   ```
   This will **compile the Lambda function** and **run Terraform** to deploy it.

---

## ⚡ How to Use

By default, the Lambda function is deployed **without triggers**. You can invoke it in two ways:

- **Scheduled trigger** (e.g., start resources in the morning, stop them in the evening).
- **Manual invocation from CLI** (e.g., start resources only when needed, ensuring unused resources remain stopped).

### **Recommended Usage**

🔹 **Better strategy**: Stop resources **every day** and only start them when needed via manual invocation. This prevents forgotten resources from staying active.

---

## ⚙️ Parameters

The Lambda function accepts several parameters to customize its behavior:

| Parameter    | Required | Default         | Description                                                                      |
| ------------ | -------- | --------------- | -------------------------------------------------------------------------------- |
| `action`     | ✅ Yes    | -               | Defines the action to perform: `start`, `stop`, or `list` (dry run).             |
| `regions`    | ❌ No     | All AWS regions | Comma-separated list of AWS regions to operate on (e.g., `eu-west-1,eu-west-2`). |
| `taggedOnly` | ❌ No     | `true`          | Whether to consider only tagged resources (`true` or `false`).                   |
| `tagKeys`    | ❌ No     | `cleeper`       | Comma-separated list of tag keys to filter resources.                            |
| `tagValues`  | ❌ No     | `true`          | Comma-separated list of tag values to filter resources.                          |

💡 **Tip:** Specifying the AWS regions you use can significantly **reduce execution time and cost**.

---

## 🛠️ Examples

### **1️⃣ List impacted resources in specific regions**

```json
{
  "action": "list",
  "regions": "eu-west-1,eu-west-2"
}
```

### **2️⃣ Stop resources with specific tags**

If you want to stop all resources with the tag **application** set to either `app1` or `secondapp` in `eu-west-1`:

```json
{
  "action": "stop",
  "regions": "eu-west-1",
  "tagKeys": "application",
  "tagValues": "app1,secondapp"
}
```

### **3️⃣ Stop all resources (ignoring tags)**

```json
{
  "action": "stop",
  "taggedOnly": "false"
}
```

This ensures that **all** resources are stopped, regardless of their tags.

---

## 📌 Running from CLI

To invoke the Lambda function via AWS CLI:

```bash
aws lambda invoke --function-name cleeper \
  --cli-binary-format raw-in-base64-out \
  --payload '{"action":"list", "regions":"eu-west-1", "tagKeys":"cleeper", "tagValues":"val2,val1"}' \
  --log-type Tail output | jq .LogResult -r | base64 -d
```

### **Example Output**

```
START RequestId: 71411120-0dd1-4dd7-9bf1-60ebb9a50889 Version: $LATEST
Working on region: eu-west-1
ASGs to suspend:
  terraform-20250226115254216000000003
EC2 instances to terminate:
  i-0c4db9d0889e307fb
EC2 instances to stop:
  i-06e32396035cf6c30
RDS Clusters to stop:
  aurora-cluster-demo
  aurora-postgres-cluster-demo
RDS instances to stop:
  terraform-20250226124539700800000001
END RequestId: 71411120-0dd1-4dd7-9bf1-60ebb9a50889
REPORT RequestId: 71411120-0dd1-4dd7-9bf1-60ebb9a50889 Duration: 592.21 ms Billed Duration: 593 ms Memory Size: 128 MB Max Memory Used: 35 MB
```

🔹 This example runs a **list action** to preview which resources would be stopped or started in `eu-west-1` using custom tags.

---

## 🎯 Conclusion

Cleeper is an **efficient tool** for managing AWS resources by ensuring they are only active when needed. By leveraging **Lambda, Terraform, and tagging mechanisms**, you can easily **reduce costs** and **automate resource management**.

🔥 **Use Cleeper to take full control over your AWS environment!** 🚀

