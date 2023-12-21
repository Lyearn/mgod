---
title: Installation
---

## Requirements

- Go 1.18 or higher
- MongoDB 3.6 and higher

:::warning
For MongoDB version **<4.4**, please create the collection in MongoDB before creating an `EntityMongoModel` using `mgod` for the same.
Refer to [this MongoDB limitations](https://www.mongodb.com/docs/manual/reference/limits/#operations) for more information.
:::

## Installation

```bash
go get github.com/Lyearn/mgod
```

As simple as that!!

Make sure that Go Mongo Driver is also installed. If not already, add it as follows -

```bash
go get go.mongodb.org/mongo-driver
```
