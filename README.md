<!-- markdownlint-disable MD033 MD013 -->
<h1 align="center">
    <br>
        OpenTelemetry Action
    <br>
</h1>
<h4 align="center">
    Upload OpenTelemetry traces of a GitHub actions workflow run
</h4>
<p align="center">
    <a href="https://github.com/MNThomson/otel-action/commits">
        <img
            src="https://img.shields.io/github/last-commit/MNThomson/otel-action?style=for-the-badge"
            alt="Last GitHub Commit"
        >
    </a>
</p>
<!-- markdownlint-enable -->

---

<!-- markdownlint-disable-next-line MD002 -->
## About

`otel-action` is currently in early Alpha.
Branch `master` is not guaranteed to be stable and breaking changes may be
introduced without notice.


## Example usage

```yaml
- name: Upload OTEL traces
  uses: MNThomson/otel-action@master
  with:
    endpoint: ${{ secrets.ENDPOINT }}
    headers: ${{ secrets.HEADERS }}
    service_name: "MyDatasetName"
```
