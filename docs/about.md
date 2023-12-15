---
title: About
---

## What is mgod?

`mgod` is a MongoDB ODM specifically designed for Go. It provides a structured way to map Go models to MongoDB collections, simplifying database interactions in Go applications.

## Why use mgod?

Creating `mgod` was driven by the need to **simplify MongoDB interactions** while keeping **one schema for all** in Go. Traditionally, working with MongoDB in Go involved either using separate structs for database and service logic or manually converting service structs to MongoDB documents, a process that was both time-consuming and error-prone. This lack of integration often led to redundant coding, especially when dealing with union types or adding meta fields for each MongoDB operation.

Inspired by the easy interface of MongoDB handling using [Mongoose](https://github.com/Automattic/mongoose) and [Typegoose](https://github.com/typegoose/typegoose) libraries available in Node, `mgod` aims to streamline these processes. It offers a more integrated approach, reducing the need to duplicate code and enhancing type safety, making MongoDB operations more intuitive and efficient in Go.
