# Paginator

This intends to be a plugin for paginating golang.

Currently only supports [gorm](https://github.com/jinzhu/gorm).


### Install

`go get -u github.com/Prabandham/paginator`

### Usage.
```
p import github.com/Prabandham/paginator


type User struct{

    gorm.Model

    Name string
    .
    .
    .
}

func FindAllUsers() []User {
    var users  []User
    db, err := gorm.Open("postgres", ....)
    order_by := []string{"name asc"}
    paginator := p.Paginator{DB: &db, OrderBy: order_by, Page: "1", PerPage: "10"}
    data := paginator.Paginate(&users)

    //data -> {
    //          total_records: 1000,
    //          records: [{Name: "Srinidhi"}, {Name: "Prabandham"},...],
    //          current_page: 1,
    //          total_pages: 100
    //        }
}

