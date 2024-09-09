import React from "react";
import style from "../common.module.css";
import { useNavigate } from "react-router-dom";
import { useState, useEffect } from "react";
import axios from "axios";

function BorrowHistory(props) {
  const navigate = useNavigate();
  const userInfo = props.userInfo;
  let [history, setHistory] = useState([]);
  const userId = userInfo.userId;
  const baseUrl = "http://localhost:5000/api/";
  const [selectedList, setSelectedList] = useState(new Set());
  const [active, setActive] = useState(false);
  useEffect(() => {
    axios
      .get(
        `${baseUrl}past-borrowed-books/booksinfo/?user_id=${userInfo.userId}`
      )
      .then((response) => {
        console.log(`Response: ${JSON.stringify(response.data.data)}`);
        const seen = new Set();
        const uniqueArray = response.data.data.filter((item) => {
          const duplicate = seen.has(item.BookId);
          seen.add(item.BookId);
          return !duplicate;
        });
        setHistory(uniqueArray);
        // const newSet = new Set();
        // uniqueArray.forEach((ele) => {
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

  const sendDeleteRequest = (lists) => {
    const newList = [];
    lists.forEach((element) => {
      newList.push({
        user_id: userId,
        book_id: element,
      });
    });
    const payload = {
      delete_data: newList,
    };
    axios
      .delete(`${baseUrl}delete-past-borrowed-books`, { data: payload })
      .then((response) => {
        setActive(!active);
        console.log(`Received response status: ${response.status}`);
      })
      .catch((err) =>
        console.log(
          "Error received during delete borrowed books request: ",
          err
        )
      );
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
              {" "}
              Lend{" "}
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
          <table className={style.borrowHistory}>
            <thead>
              <tr>
                <th scope="col1">Title</th>
                <th scope="col1">Author</th>
                <th scope="col2">Publications</th>
              </tr>
            </thead>
            <tbody>
              {history.map((ele, index) => {
                return (
                  <tr
                    key={index}
                    // style={{
                    //   backgroundColor: "#d6dcd7",
                    // }}
                    style={{
                      backgroundColor: selectedList.has(ele.BookId)
                        ? "rgb(139, 161, 87)"
                        : "#d6dcd7",
                      cursor: "pointer",
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
          <button
            className={style.returnCartButton}
            style={{
              left: "650px",
              cursor: selectedList.size > 0 ? "pointer" : "no-drop",
            }}
            disabled={selectedList.size > 0 ? false : true}
            onClick={() => sendDeleteRequest(selectedList)}
          >
            Delete
          </button>
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
                onClick={() => navigate("/", { state: {}, replace: true })}
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

export default BorrowHistory;
