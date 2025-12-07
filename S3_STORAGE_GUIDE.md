# S3-Compatible Storage Guide

## Understanding S3-Compatible Storage

### What is S3?
Amazon S3 (Simple Storage Service) is Amazon's cloud object storage service. Over time, its API became the **industry standard** for object storage.

### What does "S3-Compatible" mean?
When a service is "S3-compatible," it means it implements the same API (Application Programming Interface) as Amazon S3. This allows you to use the same tools, libraries, and code to interact with different storage providers.

### Think of it Like This:
```
┌─────────────────────────────────────────────────────┐
│            Your Go Application Code                 │
│         (Never needs to change!)                    │
└────────────────┬────────────────────────────────────┘
                 │
                 │ Uses AWS SDK (S3 Protocol)
                 │
     ┌───────────┴──────────┐
     │                      │
     ▼                      ▼
┌─────────┐         ┌──────────────┐
│ AWS S3  │   OR    │ Cloudflare R2│   OR   Any S3-compatible service
└─────────┘         └──────────────┘
```

**Analogy**: Think of S3 API like USB-C standard
- Different companies make USB-C devices (Apple, Samsung, Dell)
- All use the same USB-C connector/protocol
- One cable works with all devices

Similarly:
- AWS created S3 API standard
- Other companies implement the same API
- One SDK (AWS SDK) works with all providers

---

## Why Use AWS SDK for Non-AWS Services?

### The Power of Standardization

```go
// This is in our code (internal/services/s3_service.go)
import (
    "github.com/aws/aws-sdk-go-v2/service/s3"  // AWS SDK
)

// But we can connect to ANY S3-compatible service!
```

**Why this works:**
1. AWS SDK is just a **client library** that speaks "S3 language"
2. It doesn't care WHERE the storage is - just that it responds to S3 API calls
3. The `S3_ENDPOINT` configuration tells it where to send requests

### What the Code Does

```go
// From internal/services/s3_service.go

// This redirects AWS SDK to different storage providers
customResolver := aws.EndpointResolverWithOptionsFunc(
    func(service, region string, options ...interface{}) (aws.Endpoint, error) {
        return aws.Endpoint{
            URL: cfg.S3Endpoint,  // Your custom endpoint (Cloudflare, DigitalOcean, etc.)
        }, nil
    },
)

// AWS SDK configuration
awsCfg, _ := awsConfig.LoadDefaultConfig(context.TODO(),
    awsConfig.WithRegion(cfg.S3Region),
    awsConfig.WithEndpointResolverWithOptions(customResolver),  // Use custom endpoint
    awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
        cfg.S3Key,     // Your provider's access key
        cfg.S3Secret,  // Your provider's secret key
        "",
    )),
)
```

---

## Switching Between Storage Providers

### The Magic: Only Change `.env` File!

**No code changes needed** - just update your environment variables.

---

## 🔵 Option 1: Cloudflare R2 (Current Setup)

### Features
- ✅ No egress fees (free data transfer out)
- ✅ Cheaper than AWS S3
- ✅ Fast global CDN
- ✅ S3-compatible API
- 💰 **Cost**: $0.015/GB storage, $0 egress

### Setup Steps

1. **Create Cloudflare Account** → [dash.cloudflare.com](https://dash.cloudflare.com)

2. **Create R2 Bucket**
   - Go to R2 Object Storage
   - Click "Create bucket"
   - Enter bucket name (e.g., `ecommercegolang`)

3. **Get API Credentials**
   - Click "Manage R2 API Tokens"
   - Click "Create API Token"
   - Set permissions (Admin Read & Write)
   - Copy **Access Key ID** and **Secret Access Key**

4. **Enable Public Access**
   - Go to bucket → Settings
   - Under "Public Access", click "Allow Access"
   - Copy the **Public R2.dev URL** (e.g., `https://pub-xxxxx.r2.dev`)

5. **Configure `.env`**
```env
S3_ENDPOINT=https://your-account-id.r2.cloudflarestorage.com
S3_KEY=your-r2-access-key-id
S3_SECRET=your-r2-secret-access-key
S3_BUCKET=ecommercegolang
S3_REGION=auto
S3_PUBLIC_URL=https://pub-xxxxx.r2.dev
```

6. **Find Your Account ID**
   - R2 Dashboard → Copy the endpoint URL
   - Format: `https://[ACCOUNT_ID].r2.cloudflarestorage.com`

---

## 🟠 Option 2: Amazon S3

### Features
- ✅ Industry leader, most mature
- ✅ Extensive AWS ecosystem integration
- ✅ 99.999999999% durability (11 nines)
- ✅ Advanced features (versioning, lifecycle, etc.)
- 💰 **Cost**: $0.023/GB storage + egress fees

### Setup Steps

1. **Create AWS Account** → [aws.amazon.com](https://aws.amazon.com)

2. **Create S3 Bucket**
   - AWS Console → S3 → Create bucket
   - Choose bucket name (globally unique)
   - Select region (e.g., `us-east-1`)
   - Uncheck "Block all public access" for public images

3. **Create IAM User**
   - AWS Console → IAM → Users → Add user
   - Enable "Programmatic access"
   - Attach policy: `AmazonS3FullAccess`
   - Save **Access Key ID** and **Secret Access Key**

4. **Set Bucket Policy (for public read)**
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "PublicReadGetObject",
      "Effect": "Allow",
      "Principal": "*",
      "Action": "s3:GetObject",
      "Resource": "arn:aws:s3:::your-bucket-name/*"
    }
  ]
}
```

5. **Configure `.env`**
```env
S3_ENDPOINT=
# Leave empty - AWS SDK uses default endpoint
S3_KEY=AKIAIOSFODNN7EXAMPLE
S3_SECRET=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
S3_BUCKET=my-ecommerce-bucket
S3_REGION=us-east-1
S3_PUBLIC_URL=https://my-ecommerce-bucket.s3.us-east-1.amazonaws.com
```

**Note**: When `S3_ENDPOINT` is empty, AWS SDK automatically uses AWS S3 endpoints.

---

## 🔷 Option 3: DigitalOcean Spaces

### Features
- ✅ Simple pricing ($5/month for 250GB)
- ✅ Built-in CDN
- ✅ Easy to use
- ✅ S3-compatible API
- 💰 **Cost**: $5/month (250GB storage + 1TB transfer)

### Setup Steps

1. **Create DigitalOcean Account** → [digitalocean.com](https://digitalocean.com)

2. **Create Space**
   - Cloud → Spaces → Create Space
   - Choose region (e.g., NYC3, SFO3, AMS3)
   - Enter Space name (e.g., `ecommerce-images`)
   - Enable CDN (optional but recommended)

3. **Generate API Keys**
   - API → Spaces Keys → Generate New Key
   - Enter key name
   - Copy **Access Key** and **Secret Key**

4. **Set Public Access**
   - Go to Space settings
   - Set File Listing to "Public"

5. **Configure `.env`**
```env
S3_ENDPOINT=https://nyc3.digitaloceanspaces.com
S3_KEY=DO00ABCDEFGHIJKLMNOP
S3_SECRET=your-spaces-secret-key-here
S3_BUCKET=ecommerce-images
S3_REGION=nyc3
S3_PUBLIC_URL=https://ecommerce-images.nyc3.digitaloceanspaces.com
```

**Available Regions**: `nyc3`, `sfo3`, `ams3`, `sgp1`, `fra1`

---

## 🟣 Option 4: Backblaze B2

### Features
- ✅ Cheapest storage option
- ✅ No minimum storage
- ✅ S3-compatible API
- 💰 **Cost**: $0.005/GB storage (1/4 the price of S3)

### Setup Steps

1. **Create Backblaze Account** → [backblaze.com](https://backblaze.com/b2)

2. **Create Bucket**
   - B2 Cloud Storage → Buckets → Create Bucket
   - Enter bucket name
   - Set to "Public" for images

3. **Create Application Key**
   - App Keys → Add New Application Key
   - Select bucket or allow all
   - Copy **keyID** and **applicationKey**

4. **Get S3-Compatible Endpoint**
   - Format: `https://s3.{region}.backblazeb2.com`
   - Example: `https://s3.us-west-002.backblazeb2.com`

5. **Configure `.env`**
```env
S3_ENDPOINT=https://s3.us-west-002.backblazeb2.com
S3_KEY=your-b2-key-id
S3_SECRET=your-b2-application-key
S3_BUCKET=ecommerce-images
S3_REGION=us-west-002
S3_PUBLIC_URL=https://f002.backblazeb2.com/file/ecommerce-images
```

---

## 🟢 Option 5: MinIO (Self-Hosted)

### Features
- ✅ Run on your own server
- ✅ 100% S3-compatible
- ✅ No cloud costs
- ✅ Full control
- 💰 **Cost**: Only server hosting costs

### Setup Steps

1. **Install MinIO** (Docker example)
```bash
docker run -d \
  -p 9000:9000 \
  -p 9001:9001 \
  --name minio \
  -e "MINIO_ROOT_USER=minioadmin" \
  -e "MINIO_ROOT_PASSWORD=minioadmin" \
  -v /data/minio:/data \
  minio/minio server /data --console-address ":9001"
```

2. **Access Console**
   - Open `http://localhost:9001`
   - Login with `minioadmin` / `minioadmin`

3. **Create Bucket**
   - Buckets → Create Bucket
   - Name: `ecommerce`
   - Set access policy to "Public"

4. **Create Access Keys**
   - Identity → Service Accounts → Create Service Account
   - Save Access Key and Secret Key

5. **Configure `.env`**
```env
S3_ENDPOINT=http://localhost:9000
S3_KEY=minioadmin
S3_SECRET=minioadmin
S3_BUCKET=ecommerce
S3_REGION=us-east-1
S3_PUBLIC_URL=http://localhost:9000/ecommerce
```

---

## 🔴 Option 6: Wasabi

### Features
- ✅ Hot cloud storage
- ✅ 80% cheaper than AWS S3
- ✅ Free egress (no bandwidth fees)
- ✅ S3-compatible API
- 💰 **Cost**: $5.99/TB/month, $0 egress

### Setup Steps

1. **Create Wasabi Account** → [wasabi.com](https://wasabi.com)

2. **Create Bucket**
   - Buckets → Create Bucket
   - Choose region
   - Enable "Everyone" read access for public images

3. **Create Access Keys**
   - Access Keys → Create New Access Key
   - Copy **Access Key** and **Secret Key**

4. **Get Endpoint URL**
   - Format: `https://s3.{region}.wasabisys.com`
   - Regions: `us-east-1`, `us-east-2`, `us-west-1`, `eu-central-1`

5. **Configure `.env`**
```env
S3_ENDPOINT=https://s3.us-east-1.wasabisys.com
S3_KEY=your-wasabi-access-key
S3_SECRET=your-wasabi-secret-key
S3_BUCKET=ecommerce-images
S3_REGION=us-east-1
S3_PUBLIC_URL=https://s3.us-east-1.wasabisys.com/ecommerce-images
```

---

## 📊 Cost Comparison

| Provider | Storage Cost/GB/Month | Egress/GB | Free Tier | Best For |
|----------|----------------------|-----------|-----------|----------|
| **Cloudflare R2** | $0.015 | **$0** | 10GB free | High traffic apps |
| **AWS S3** | $0.023 | $0.09 | 5GB/12mo | AWS ecosystem |
| **DigitalOcean** | $0.02 | Included | None | Simple pricing |
| **Backblaze B2** | **$0.005** | $0.01 | 10GB free | Cost-sensitive |
| **Wasabi** | $0.0059 | **$0** | None | Large storage |
| **MinIO** | Server cost | $0 | N/A | Self-hosted |

---

## 🔄 How to Switch Providers

### Step-by-Step Migration

1. **Set up new provider** (follow steps above)
2. **Update `.env` file** with new credentials
3. **Restart your application**
```bash
go run ./cmd/api/main.go
```
4. **Test upload**
```bash
curl -X POST http://localhost:8080/api/upload/product \
  -H "Authorization: Bearer <admin-token>" \
  -F "image=@test.jpg"
```

5. **(Optional) Migrate existing files**
   - Use tools like `rclone` or `s3cmd` to copy files between providers
   - Or write a migration script

### Example Migration Script
```bash
# Install rclone
# Configure source and destination

rclone copy cloudflare-r2:ecommercegolang aws-s3:my-bucket
```

---

## 🛠️ Code Configuration Breakdown

### How Our Code Handles Different Providers

```go
// internal/services/s3_service.go

func InitS3() {
    cfg := appConfig.Cfg
    
    // Custom endpoint resolver
    // This is what makes switching possible!
    customResolver := aws.EndpointResolverWithOptionsFunc(
        func(service, region string, options ...interface{}) (aws.Endpoint, error) {
            return aws.Endpoint{
                URL: cfg.S3Endpoint,  // From .env file
            }, nil
        },
    )
    
    awsCfg, _ := awsConfig.LoadDefaultConfig(context.TODO(),
        awsConfig.WithRegion(cfg.S3Region),           // From .env
        awsConfig.WithEndpointResolverWithOptions(customResolver),
        awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
            cfg.S3Key,     // From .env
            cfg.S3Secret,  // From .env
            "",
        )),
    )
    
    s3Client = s3.NewFromConfig(awsCfg)
}
```

### Environment Variables Explained

| Variable | Purpose | Example |
|----------|---------|---------|
| `S3_ENDPOINT` | Where to send requests | `https://xxx.r2.cloudflarestorage.com` |
| `S3_KEY` | Authentication - Access Key | `AKIAIOSFODNN7EXAMPLE` |
| `S3_SECRET` | Authentication - Secret Key | `wJalrXUtnFEMI/K7MDENG...` |
| `S3_BUCKET` | Which bucket to use | `ecommerce-images` |
| `S3_REGION` | Region/datacenter | `us-east-1`, `auto`, `nyc3` |
| `S3_PUBLIC_URL` | Public URL for uploaded files | `https://pub-xxx.r2.dev` |

---

## ✅ Testing Your Setup

### 1. Check Configuration
```bash
# View loaded config in terminal when server starts
go run ./cmd/api/main.go
```

### 2. Test Upload
```bash
curl -X POST http://localhost:8080/api/upload/product \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -F "image=@/path/to/test.jpg"
```

### 3. Verify Response
```json
{
  "url": "https://pub-xxxxx.r2.dev/1234567890.jpg"
}
```

### 4. Check in Browser
- Open the returned URL
- Image should display (not download)
- Check browser network tab for correct content-type

---

## 🐛 Troubleshooting

### Issue: "InvalidAccessKeyId"
**Solution**: Double-check `S3_KEY` in `.env`

### Issue: "SignatureDoesNotMatch"
**Solution**: Verify `S3_SECRET` is correct

### Issue: "NoSuchBucket"
**Solution**: Ensure bucket name in `S3_BUCKET` matches actual bucket

### Issue: "Connection refused"
**Solution**: Check `S3_ENDPOINT` URL is correct and accessible

### Issue: Image downloads instead of displaying
**Solution**: Check `ContentType` is set in upload code (already fixed in our code)

### Issue: "403 Forbidden"
**Solution**: 
- Verify bucket is set to public
- Check IAM/API key permissions

---

## 🔒 Security Best Practices

1. **Never commit `.env` to Git**
   - Add `.env` to `.gitignore`

2. **Use environment-specific credentials**
   - Development: Separate bucket/keys
   - Production: Different bucket/keys

3. **Limit IAM permissions**
   - Only grant S3 access, not full AWS access
   - Use bucket-specific policies

4. **Rotate credentials regularly**
   - Change keys every 90 days

5. **Use HTTPS endpoints**
   - Never use HTTP in production

6. **Enable versioning**
   - Protect against accidental deletions

---

## 📚 Additional Resources

### Official Documentation
- [AWS S3 API Reference](https://docs.aws.amazon.com/s3/index.html)
- [Cloudflare R2 Docs](https://developers.cloudflare.com/r2/)
- [DigitalOcean Spaces](https://docs.digitalocean.com/products/spaces/)
- [Backblaze B2](https://www.backblaze.com/b2/docs/)
- [MinIO Documentation](https://min.io/docs/minio/linux/index.html)
- [Wasabi Documentation](https://wasabi.com/help/)

### Tools
- [AWS CLI](https://aws.amazon.com/cli/) - Command-line S3 management
- [s3cmd](https://s3tools.org/s3cmd) - S3 client for command line
- [Cyberduck](https://cyberduck.io/) - GUI for S3 storage
- [rclone](https://rclone.org/) - Sync files between clouds

---

## 🎯 Summary

### Key Takeaways

1. **S3 API is a standard** - Not just AWS
2. **AWS SDK is just a client library** - Works with any S3-compatible service
3. **Switching is easy** - Only change `.env` file
4. **No vendor lock-in** - Move between providers freely
5. **Same code, different storage** - Flexibility without rewrites

### Benefits of This Architecture

✅ **Flexibility** - Switch providers anytime
✅ **Cost optimization** - Choose cheapest option
✅ **No code changes** - Just configuration
✅ **Risk mitigation** - Not locked to one vendor
✅ **Testing** - Use MinIO locally, cloud in production

---

**Happy Cloud Storage! ☁️📦**
