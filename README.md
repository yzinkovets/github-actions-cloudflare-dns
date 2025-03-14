# github-actions-cloudflare-dns
Github Action to work with CloudFlare DNS

Please note, at the moment it works with `CNAME` records only.

### Usage

Github Action:
```yaml
jobs:
  update-cloudflare:
    runs-on: ubuntu-latest
    steps:
      - name: Create/Update CNAME Record in Cloudflare
        uses: yzinkovets/github-actions-cloudflare-dns@v1
        with:
          cloudflare_api_token: ${{ secrets.CLOUDFLARE_API_TOKEN }}
          domain: "sub.example.com"
          target: "domain.target.com"
```