import React, { useEffect } from "react";
import style from "../common.module.css";
import { useNavigate } from "react-router-dom";
import axios from "axios";
import { useState } from "react";

function LendCart(props) {
  const navigate = useNavigate();
  const userInfo = props.userInfo;
  let [data, setData] = useState([]);
  const baseUrl = "http://localhost:5000/api/";
  const [selectedList, setSelectedList] = useState(new Set());
  const [active, setActive] = useState(false);
  useEffect(() => {
    axios
      .get(`${baseUrl}add-to-cart-books/booksinfo?user_id=${userInfo.userId}`, {
        headers: { "Content-Type": "application/json" },
      })
      .then((response) => {
        console.log(`Response: ${JSON.stringify(response.data.data)}`);
        const res = response.data.data;
        setData(res);
        // const newSet = new Set();
        // res.forEach((ele) => {
        //   if (selectedList.has(ele.BookId)) {
        //     newSet.add(ele.BookId);
        //   }
        // });
        // setSelectedList(newSet);
        setSelectedList(new Set());
      })
      .catch((error) => {
        console.log(error);
      });
  }, [active, userInfo.userId]);
  const sendLendRequest = (list) => {
    const newList = [];
    list.forEach((ele) => {
      const obj = {
        user_id: userInfo.userId,
        book_id: ele,
      };
      newList.push(obj);
    });
    console.log(newList);
    const payload = {
      lend_data: newList,
    };
    // axios
    //   .post(
    //     `${baseUrl}borrow-books`,
    //     payload
    //     // {
    //     //   headers: {
    //     //     "content-type": "application/json",
    //     //   },
    //     //   responseType: "json",
    //     //   withCredentials: true,
    //     // }
    //   )
    //   .then((response) => {
    //     console.log(response.status, response.data);
    //     setActive(!active);
    //   });
    axios
      .put(`${baseUrl}borrow-books`, payload)
      .then((response) => {
        setActive(!active);
        console.log(response);
      })
      .catch((err) =>
        console.log("Encountered error during lend requeest: ", err)
      );
  };
  const sendDeleteRequest = (lists) => {
    const newList = [];
    lists.forEach((ele) => {
      newList.push({
        user_id: userInfo.userId,
        book_id: ele,
      });
    });
    const deletePayload = {
      delete_data: newList,
    };
    axios
      .delete(`${baseUrl}delete-add-to-cart-books`, { data: deletePayload })
      .then((response) => {
        setActive(!active);
        console.log("Borrow Books delete request status: ", response.status);
      })
      .catch((err) => console.log(`Error occurred in delete request: ${err}`));
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
                  state: { isLoggedIn: props.isLoggedIn, userInfo: userInfo },
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
          </ul>
        </section>
        <section>
          <table className={style.lendCartCentre}>
            <thead>
              <tr>
                <th scope="col1">Title</th>
                <th scope="col1">Author</th>
                <th scope="col2">Publications</th>
              </tr>
            </thead>
            <tbody>
              {data.map((ele) => {
                return (
                  <tr
                    key={ele.BookId}
                    style={{
                      backgroundColor: selectedList.has(ele.BookId)
                        ? "rgb(139, 161, 87)"
                        : "#d6dcd7",
                    }}
                    onClick={() => {
                      setSelectedList((prevSelected) => {
                        const newSelected = new Set(prevSelected);
                        if (newSelected.has(ele.BookId)) {
                          newSelected.delete(ele.BookId);
                        } else {
                          newSelected.add(ele.BookId);
                        }
                        return newSelected;
                      });
                    }}
                  >
                    <td>{ele.Title}</td>
                    <td>{ele.Author}</td>
                    <td>{ele.Publications}</td>
                  </tr>
                );
              })}
            </tbody>
          </table>
          <div className={style.buttonDiv}>
            <button
              className={style.lendCartButton}
              style={{ cursor: selectedList.size > 0 ? "pointer" : "no-drop" }}
              disabled={selectedList.size > 0 ? false : true}
              onClick={() => sendLendRequest(selectedList)}
            >
              Lend Request
            </button>
            <button
              className={style.lendCartButton}
              style={{ cursor: selectedList.size > 0 ? "pointer" : "no-drop" }}
              disabled={selectedList.size > 0 ? false : true}
              onClick={() => sendDeleteRequest(selectedList)}
            >
              Delete
            </button>
          </div>
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
                onClick={() =>
                  navigate("/", {
                    replace: true,
                    state: {},
                  })
                }
              >
                Logout
              </button>
            </ul>
          </div>
        </section>
      </div>
    </div>
  );
}

export default LendCart;
