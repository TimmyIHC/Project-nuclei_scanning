## Research Initiative: Vulnerability Scanning with Nuclei in the Chaos-Bugbounty Suite

As a component of our ongoing research efforts, this project conducts exhaustive scans for vulnerabilities within the chaos-bugbounty suite utilizing Nuclei—a cutting-edge, open-source scanner. This comprehensive approach generates in-depth reports for each assessed domain. Subsequently, these reports are securely transmitted to a designated Amazon S3 storage bucket.

The adoption of cloud services for this project was strategic; vulnerability scanning can be an extensive process, sometimes extending over several days. Leveraging cloud infrastructure ensures the scans can run uninterrupted without occupying local resources.

### Prerequisites for Use:
Before you begin, ensure the following steps are completed:

1. Set up an Amazon Web Services (AWS) account.
2. Provision an Amazon EC2 instance to deploy the code.
3. Install Docker onto the EC2 instance.
4. Modify the code at line number 114 to reflect the AWS region where your S3 bucket was created.
5. At line number 77, update the bucket name to the one into which you intend to upload the reports.

### Execution Instructions:
With Docker and the respective Go file situated on your EC2 machine, execute the commands below to build and run your Docker container:

```sh
docker build -t your_image_name .
docker run your_image_name
```

Ensure that you replace `your_image_name` with the name you have assigned to your Docker image.
