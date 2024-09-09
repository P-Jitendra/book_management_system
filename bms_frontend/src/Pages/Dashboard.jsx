import React from "react";
import style from "../common.module.css";
import { useNavigate } from "react-router-dom";
import axios from "axios";
import { useState, useEffect } from "react";

export default function Dashboard(props) {
  const userInfo = props.userInfo;
  const navigate = useNavigate();
  const userId = userInfo.userId;
  // ("2017A1PS1030H");
  const baseUrl = "http://localhost:5000/api/";
  const [topBooksList, setTopBooksList] = useState([]);
  const n = 5;
  useEffect(() => {
    axios
      .get(`${baseUrl}booksinfo/top-n-books/${n}`)
      .then((response) => {
        console.log(`Top 5 books data:${JSON.stringify(response.data.data)}`);
        setTopBooksList(response.data.data);
      })
      .catch((err) => {
        console.log(`Error occurred while fetch top 5 books: ${err}`);
      });
  }, []);

  const createAddToCartReq = async (params) => {
    const arr = [
      {
        book_id: params.BookId,
        user_id: userId,
      },
    ];
    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    // const response = await fetch(`${baseUrl}add-to-cart-books/`, {
    //   method: "POST",
    //   mode: "no-cors",
    //   headers: myHeaders,
    //   body: JSON.stringify({
    //     add_to_cart_data: arr,
    //   }),
    // });
    // console.log("Printing status: ", response.status);
    const response = await axios.post(`${baseUrl}add-to-cart-books`, {
      add_to_cart_data: arr,
    });
    console.log("Printing status: ", response.status);
  };

  // const createAddToCartReq = (ele) => {
  //   const arr = [
  //     {
  //       book_id: ele.BookId,
  //       user_id: userId,
  //     },
  //   ];
  //   const myHeaders = new Headers();
  //   myHeaders.append("Content-Type", "application/json");
  //   myHeaders.append("Aceess-Control-Allow-Origin", "*");
  //   fetch(`${baseUrl}add-to-cart-books/`, {
  //     method: "POST",
  //     mode: "no-cors",
  //     headers: myHeaders,
  //     body: JSON.stringify({
  //       add_to_cart_data: arr,
  //     }),
  //   })
  //     .then((response) => {
  //       console.log(response);
  //       console.log(response.json());
  //       response.text();
  //     })
  //     .then((data) => console.log("Printing data: ", data))
  //     // .then((res1) => {
  //     //   console.log(res1);
  //     //   console.log(res1.body);
  //     // })
  //     .catch((err) => {
  //       console.log(
  //         `Error occurred while creating add to cart request: ${err}`
  //       );
  //     });
  // };

  const handleGenreData = (e, genreId) => {
    e.preventDefault();
    navigate("/genre/" + genreId, {
      state: {
        isLoggedIn: props.isLoggedIn,
        userInfo: props.userInfo,
      },
      replace: true,
    });
  };
  return (
    <div className={style.commonStyle}>
      <section className={style.layout}>
        <h1 className={style.title}>Book Management System</h1>
      </section>
      <div className={style.layout3}>
        <section className={style.layout1}>
          <ul>
            <button
              onClick={() =>
                navigate("/dashboard", {
                  state: { isLoggedIn: true, userInfo: userInfo },
                  replace: true,
                })
              }
            >
              Home
            </button>
            <button
              onClick={() =>
                navigate("/lendCart", {
                  state: { isLoggedIn: props.isLoggedIn, userInfo: userInfo },
                  replace: true,
                })
              }
            >
              Lend
            </button>
            <button
              onClick={() =>
                navigate("/returnCart", {
                  state: { isLoggedIn: props.isLoggedIn, userInfo: userInfo },
                  replace: true,
                })
              }
            >
              Return
            </button>
            <button
              onClick={() =>
                navigate("/borrowHistory", {
                  state: { isLoggedIn: props.isLoggedIn, userInfo: userInfo },
                  replace: true,
                })
              }
            >
              Borrow History
            </button>
            {/* <li>
              <span>Lend</span>
            </li>
            <li>
              <span>Return</span>
            </li>
            <li>
              <span>Borrow History</span>
            </li> */}
          </ul>
        </section>
        <section className={style.centreElement}>
          <table>
            <tbody>
              <tr>
                <th>
                  <span>GENRES</span>
                </th>
              </tr>
              <tr key="first">
                <td>
                  <button onClick={(e) => handleGenreData(e, "Fantasy")}>
                    Fantasy
                  </button>
                  <button onClick={(e) => handleGenreData(e, "Horror")}>
                    Horror
                  </button>
                  <button
                    onClick={(e) => handleGenreData(e, "Science Fiction")}
                  >
                    Science Fiction
                  </button>
                  <button onClick={(e) => handleGenreData(e, "Action")}>
                    Action
                  </button>
                  <button onClick={(e) => handleGenreData(e, "Mystery")}>
                    Mystery
                  </button>
                </td>
              </tr>
              <tr key="second">
                <td>
                  <button onClick={(e) => handleGenreData(e, "Historical")}>
                    History
                  </button>
                  <button onClick={(e) => handleGenreData(e, "Romance")}>
                    Romance
                  </button>
                  <button onClick={(e) => handleGenreData(e, "Political")}>
                    Political
                  </button>
                  <button onClick={(e) => handleGenreData(e, "Thrillers")}>
                    Thrillers
                  </button>
                  <button onClick={(e) => handleGenreData(e, "Adventure")}>
                    Adventure
                  </button>
                </td>
              </tr>
              <tr>
                <th>
                  <span>TOP 5 Famous Books</span>
                </th>
              </tr>
              <tr className={style.topBooksTR} key={3}>
                {topBooksList.map((ele, index) => {
                  return (
                    <td className={style.topBooks} key={ele.BookId}>
                      <span>Book{index + 1}</span>
                      <button
                        key={ele.BookId}
                        onClick={() => createAddToCartReq(ele)}
                      >
                        Add to Cart
                      </button>
                    </td>
                  );
                })}
                {/* <div style={{ display: "flex" }}>
                <button
                  key={1}
                  type="submit"
                  // onClick={() => createAddToCartReq(ele)}
                  onClick={() => console.log(`Clicked now!`)}
                  disabled={false}
                ></button>
                <button
                  key={2}
                  type="submit"
                  // onClick={() => createAddToCartReq(ele)}
                  onClick={() => console.log(`Clicked now!`)}
                  disabled={false}
                ></button>
              </div> */}
                {/* <td className={style.topBooks}>
                <span>Book1</span>
                <button onClick={() => createAddToCartReq()}>Add to Cart</button>
              </td>
              <td className={style.topBooks}>
                <span>Book2</span>
                <button>Add to Cart</button>
              </td>
              <td className={style.topBooks}>
                <span>Book3</span>
                <button>Add to Cart</button>
              </td>
              <td className={style.topBooks}>
                <span>Book4</span>
                <button>Add to Cart</button>
              </td>
              <td className={style.topBooks}>
                <span>Book5</span>
                <button>Add to Cart</button>
              </td> */}
              </tr>
            </tbody>
          </table>
        </section>
        <section className={style.layout2}>
          <div className={style.container}>
            <ul>
              <button
                onClick={() =>
                  navigate("/account", {
                    state: { isLoggedIn: props.isLoggedIn, userInfo: userInfo },
                    replace: true,
                  })
                }
              >
                Account
              </button>
              <button
                onClick={() =>
                  navigate("/notifications", {
                    state: { isLoggedIn: props.isLoggedIn, userInfo: userInfo },
                    replace: true,
                  })
                }
              >
                Notifications
              </button>
              <button
                onClick={() => {
                  navigate("/", {
                    replace: true,
                    state: {},
                  });
                }}
              >
                Logout
              </button>
              {/* <li>
                <span>Account</span>
              </li>
              <li>
                <span>Notifications</span>
              </li>
              <li>
                <span>Logout</span>
              </li> */}
            </ul>
          </div>
        </section>
      </div>
    </div>
  );
}
