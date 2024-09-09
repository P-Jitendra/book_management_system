import React from "react";
import style from "../common.module.css";
import { useNavigate } from "react-router-dom";
import { useState } from "react";
import { HiOutlineEye, HiOutlineEyeOff } from "react-icons/hi";
import axios from "axios";

function Account(props) {
  const navigate = useNavigate();
  const userInfo = props.userInfo;
  const [name, setName] = useState("");
  const [userId, setUserId] = useState("");
  const [email, setEmail] = useState("");
  const [contactNo, setContactNo] = useState("");
  const [password, setPassword] = useState("");
  const [needUpdate, setNeedUpdate] = useState(false);
  const baseUrl = "http://localhost:5000/api/";
  const [visible, setVisibility] = useState(false);
  const [error, setError] = useState("");
  const [successResponse, setSuccessResponse] = useState("");
  const sendUpdadteReq = (updateObj) => {
    updateObj.contact_no = parseInt(updateObj.contact_no);
    console.log("Printing update obj: ", updateObj);
    axios
      .put(`${baseUrl}update-user`, updateObj, {
        headers: { "Content-Type": "application/json" },
      })
      .then((response) => {
        console.log(response);
        if (response.data.status === "success") {
          setSuccessResponse("User details updated successfully!");
          setError("");
          setNeedUpdate(false);
          setName("");
          setUserId("");
          setEmail("");
          setContactNo("");
          setPassword("");
        } else {
          console.log("Printing response data: ", response.data);
          setError("Updated failed!, Plz enter proper data.");
        }
      })
      .catch((err) => {
        console.log(err);
        setError("Error occurred while updating the data!");
      });
  };

  const checkAndUpdateErrAndSuccessState = () => {
    if (error.length > 0) {
      setError("");
    }
    if (successResponse.length > 0) {
      setSuccessResponse("");
    }
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
          <form className={style.AccountCentre}>
            <span className={style.AccountSpan}>Profile</span>
            <input
              type="text"
              id="Name"
              placeholder="Name"
              autoComplete="off"
              onChange={(e) => {
                checkAndUpdateErrAndSuccessState();
                setName(e.target.value);
              }}
              value={name}
              className={style.AccountInput}
            />
            <input
              type="text"
              id="UserID"
              placeholder="UserID"
              onChange={(e) => {
                checkAndUpdateErrAndSuccessState();
                setUserId(e.target.value);
              }}
              value={userId}
              className={style.AccountInput}
            />
            <input
              type="email"
              id="Email"
              placeholder="Email"
              autoComplete="off"
              value={email}
              onChange={(e) => {
                checkAndUpdateErrAndSuccessState();
                setEmail(e.target.value);
              }}
              className={style.AccountInput}
            />
            <input
              type="tel"
              id="ContactNo"
              placeholder="Contact No"
              onChange={(e) => {
                checkAndUpdateErrAndSuccessState();
                setContactNo(e.target.value);
              }}
              value={contactNo}
              className={style.AccountInput}
            />
            <input
              type={visible ? "text" : "password"}
              id="Password"
              placeholder="Password"
              autoComplete="off"
              value={password}
              onChange={(e) => {
                checkAndUpdateErrAndSuccessState();
                setPassword(e.target.value);
              }}
              className={style.AccountInput}
            />
            <span className={style.AccountPwdIcon}>
              {visible ? (
                <HiOutlineEye
                  size={"22"}
                  onClick={() => setVisibility(!visible)}
                />
              ) : (
                <HiOutlineEyeOff
                  size={"22"}
                  onClick={() => setVisibility(!visible)}
                />
              )}
            </span>
            <span className={style.AccountErrorSpan}>{error}</span>
            <span className={style.AccountSuccessSpan}>{successResponse}</span>
            <div>
              <button
                style={{ cursor: !needUpdate ? "pointer" : "no-drop" }}
                disabled={!needUpdate ? false : true}
                className={style.AccountBtn}
                onClick={(e) => {
                  e.preventDefault();
                  setNeedUpdate(!needUpdate);
                  setSuccessResponse("");
                  setError("");
                }}
              >
                Edit
              </button>
              <button
                style={{ cursor: needUpdate ? "pointer" : "no-drop" }}
                disabled={needUpdate ? false : true}
                className={style.AccountBtn}
                onClick={(e) => {
                  e.preventDefault();
                  console.log("Printing props userId: ", userInfo.userId);
                  if (userId.length === 0) {
                    setError("UserID shouldn't be empty!");
                  } else if (userId !== userInfo.userId) {
                    setError("Logged UserID should match!");
                    setUserId("");
                  } else if (!email.includes("@")) {
                    setError("Email doesn't contain @!");
                    setEmail("");
                  } else if (contactNo.length !== 10) {
                    setError("ContactNo should contain max 10 digits!");
                    setContactNo("");
                  } else if (contactNo.charAt(0) === "0") {
                    setError("ContactNo doesn't start with zero!");
                    setContactNo("");
                  } else if (isNaN(Number(contactNo))) {
                    setError("ContactNo should contain only digits!");
                  } else {
                    sendUpdadteReq({
                      name: name,
                      email: email,
                      user_id: userId,
                      contact_no: contactNo,
                      password: password,
                      category: "user",
                    });
                  }
                }}
              >
                Update
              </button>
            </div>
          </form>
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

export default Account;
