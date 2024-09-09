import React, { useState } from "react";
import style from "../common.module.css";
import { Link, useNavigate } from "react-router-dom";
import axios from "axios";
import { HiOutlineEye, HiOutlineEyeOff } from "react-icons/hi";

function Login() {
  const navigate = useNavigate();
  const [userID, setUserID] = useState("");
  const [password, setPassword] = useState("");
  const [visible, setVisibility] = useState(false);
  const [error, setError] = useState("");
  const baseUrl = "http://localhost:5000/api/";
  const sendLoginRequest = (e, currUserId, currPassword) => {
    e.preventDefault();
    if (currUserId.length > 0) {
      axios
        .post(
          `${baseUrl}check-user-existence/user-id`,
          {
            user_id: currUserId,
            password: currPassword,
          },
          {
            headers: { "Content-Type": "application/json" },
          }
        )
        .then((response) => {
          console.log(response);
          if (response.data.status === "success") {
            navigate("/dashboard", {
              state: {
                isLoggedIn: true,
                userInfo: {
                  userId: currUserId,
                },
              },
              replace: true,
            });
            setUserID("");
            setPassword("");
          } else {
            navigate("/");
          }
        })
        .catch((err) => {
          console.log(err);
          if (err.response.data.msg === "Password doesn't match!") {
            setError("Invalid Password!");
          } else if (err.response.data.msg === "UserId doesn't exist!") {
            setError("User doesn't exist!");
          }
        });
    } else {
      console.log("Email/UserId field is empty so cannot send get request.");
    }
  };
  return (
    <div className={style.Login}>
      <form className={style.LoginCentre}>
        <span className={style.LoginSpan}>BMS App</span>
        <input
          id="userId"
          type="text"
          placeholder="UserId"
          onChange={(e) => {
            setUserID(e.target.value);
          }}
          className={style.LoginInput}
          autoComplete="off"
        />
        <input
          id="password"
          type={visible ? "text" : "password"}
          placeholder="Password"
          onChange={(e) => setPassword(e.target.value)}
          value={password}
          className={style.LoginInput}
          autoComplete="off"
        ></input>
        <span className={style.LoginPwdIcon}>
          {visible ? (
            <HiOutlineEye size={"22"} onClick={() => setVisibility(!visible)} />
          ) : (
            <HiOutlineEyeOff
              size={"22"}
              onClick={() => setVisibility(!visible)}
            />
          )}
        </span>
        <span className={style.ErrorSpan}>{error}</span>
        <button
          className={style.LoginBtn}
          onClick={(e) => sendLoginRequest(e, userID, password)}
        >
          Login
        </button>

        <span>
          Don't have an account, <Link to="/register">Register</Link>
        </span>
      </form>
    </div>
  );
}

export default Login;
