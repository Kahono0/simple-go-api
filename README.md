# simple-go-api

This is a simple graphql api written in go. It uses the [gqlgen](https://gqlgen.com/) library to generate the graphql server.
It allows users to create, read, update, and delete items and orders.

## Authentication

The app uses oidc for authentication. It uses the [go-oidc]("github.com/coreos/go-oidc") library to authenticate users. The app is configured to use the Google Identity Platform for authentication.
Once a user clicks on the `login with google` button, they are redirected to the google login page where they can sign in with their google account. Once signed in, a cookie is set which is used to authenticate the user for subsequent requests.

## Database

The app uses a postgres database which is an instance on render.

## Deployment and CI/CD

The app uses the [render](https://render.com/) platform for deployment.

## Using the app

Visit the app at [https://simple-go-api.onrender.com/](https://simple-go-api.onrender.com/)
You will be asked to sign in with your google account. Once signed in, you can use the graphql playground to interact with the api.

Below are some sample queries and mutations you can use to interact with the api.

```graphql
# your profile
query {
  me {
    id
    email
    name
  }
}

# create an item
mutation {
  createItem(
    name: "Item 2"
    description: "This is the second item"
    price: 100.00
  ) {
    id
    name
    description
  }
}

# get all items
query {
  items {
    id
    name
    description
    price
  }
}

# create an order
mutation {
  createOrder(
    input: {
      contact: "<phone number>"
      items: [
        "5a31da82-e844-43fc-8ac0-b69aab50cbfe"
        "2afe6cba-e550-48ee-b8fd-32be2b336104"
      ]
    }
  ) {
    id
    items {
      id
      name
      description
      price
    }
    total
    status
    createdAt
  }
}

# get all orders you created
query {
  orders {
    id
    items {
      id
      name
      description
      price
    }
    total
    status
    createdAt
  }
}
```

# How to run the app locally

Clone the repository and navigate to the root directory of the project. Run the following commands to start the app.

```bash
go mod tidy
go run main.go
```

You will need to create a `.env` file in the root directory of the project according to the `.env.example` file.

The app will start on port 8080. You can visit the graphql playground at [http://localhost:8080/](http://localhost:8080/)

## Testing

To run the tests, run the following command in the root directory of the project.

```bash
go test ./graph/graph_test.go
```
