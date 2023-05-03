# `nip`: Easy Frontend Containers

`nip` is a simple, fast container image builder for frontend applications.

It's ideal for use cases where your project outputs static files that just need to be served (e.g., a React app).

`nip` builds images by effectively executing `npm ci && npm run build` on your local machine, and as such doesn't require `docker` to be installed, not does it require `root` permissions.
This can make it a good fit for lightweight CI/CD use cases.

## Install `nip` and get started!

### Acknowledgements

This work is inspired by [`ko`](https://github.com/ko-build/ko).
