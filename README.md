# Golang GIN and PGX Test example

<p align="center">
  <img src="screen.png" />
</p>

This is simple [Golang](https://golang.org) test example using [Gin](1), and [pgx](2) driver for [PostgreSQL](https://postgresql.org). This simple example project also was a reminder for me in case I need to do same test using similar stack in the future (as for few days ago, i am really stuck testing my project and dig out my brain searching here and there for completing the testing task). Hopefully it will help you too.


#### Package used in this test project:
- [x] [Gin](1)
- [x] [Pgx](2)
- [x] [Testify](3)
- [x] [PgxMock](4)


### Runing Test

Open your terminal `bash, zsh, etc`

#### This will run test via terminal

```bash
make test
```
<p align="center">
  <img src="test.png" />
</p>


#### This will show the test with cover in browser

```bash
make test-cover
```
<p align="center">
  <img src="test-browser.png" />
</p>



## LICENSE
[MIT](https://github.com/reshimahendra/gin-starter/blob/main/LICENSE)

[1]:https://github.com/gin-gonic/gin
[2]:https://github.com/jackc/pgx/v4
[3]:https://github.com/stretcr/testify
[4]:https://github.com/pashagolub/pgxmock
