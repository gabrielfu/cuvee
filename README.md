# cuvee

A simple web application for managing wine cellars.

<b>View all your wines and purchases</b>

![](assets/cellar.png)

<b>Add a new wine</b>

Click "Search for an image online" to add a picture of the wine.

![](assets/add_wine.png)

<b>Search by name, vintage, country or region</b>

![](assets/search.png)

<b>View the vintage chart of a given region</b>

![](assets/with_rating.png)

<b>Ask AI to suggest a suitable vintage chart region</b>

![](assets/add_region.png)
![](assets/suggested_region.png)

## Run

#### Backend

It is recommended to run with [air](https://github.com/air-verse/air).

```shell
cd backend
go mod download
air
```

#### Frontend

```shell
cd frontend
npm i
npm run dev
```
