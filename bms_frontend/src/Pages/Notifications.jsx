import React, { useState } from "react";
import style from "../common.module.css";
import { useNavigate } from "react-router-dom";

function Notifications(props) {
  const navigate = useNavigate();
  const [notifications, setNotifications] = useState([]);
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
                  state: {
                    isLoggedIn: props.isLoggedIn,
                    userInfo: props.userInfo,
                  },
                  replace: true,
                })
              }
            >
              Home
            </button>
            <button
              onClick={() =>
                navigate("/lendCart", {
                  state: {
                    isLoggedIn: props.isLoggedIn,
                    userInfo: props.userInfo,
                  },
                  replace: true,
                })
              }
            >
              Lend
            </button>
            <button
              onClick={() =>
                navigate("/returnCart", {
                  state: {
                    isLoggedIn: props.isLoggedIn,
                    userInfo: props.userInfo,
                  },
                  replace: true,
                })
              }
            >
              Return
            </button>
            <button
              onClick={() =>
                navigate("/borrowHistory", {
                  state: {
                    isLoggedIn: props.isLoggedIn,
                    userInfo: props.userInfo,
                  },
                  replace: true,
                })
              }
            >
              Borrow History
            </button>
          </ul>
        </section>
        {notifications.length > 0 ? (
          <h1>New Notifications</h1>
        ) : (
          <span className={style.notificationsSpan}>No Notifications</span>
        )}
        <section className={style.layout2}>
          <div className={style.container}>
            <ul>
              <button
                onClick={() =>
                  navigate("/account", {
                    state: {
                      isLoggedIn: props.isLoggedIn,
                      userInfo: props.userInfo,
                    },
                    replace: true,
                  })
                }
              >
                Account
              </button>
              <button
                onClick={() => {
                  setNotifications([]);
                  navigate("/notifications", {
                    state: {
                      isLoggedIn: props.isLoggedIn,
                      userInfo: props.userInfo,
                    },
                    replace: true,
                  });
                }}
              >
                Notifications
              </button>
              <button
                onClick={() =>
                  navigate("/", {
                    state: {},
                    replace: true,
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

export default Notifications;
