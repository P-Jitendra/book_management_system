import React, { useState } from "react";
import style from "../common.module.css";
import { Link, useNavigate } from "react-router-dom";
import axios from "axios";
import { HiOutlineEye, HiOutlineEyeOff } from "react-icons/hi";

function Register() {
  const [name, setName] = useState("");
  const [userId, setUserId] = useState("");
  const [email, setEmail] = useState("");
  const [contactNo, setContactNo] = useState("");
  const [password, setPassword] = useState("");
  const baseUrl = "http://localhost:5000/api/";
  const navigate = useNavigate();
  const [visible, setVisibility] = useState(false);
  const [error, setError] = useState("");

  const sendRegisterReq = (registerObj) => {
    registerObj.ContactNo = parseInt(registerObj.ContactNo);
    const payload = {
      user_data: [registerObj],
    };
    console.log("Payload: ", payload);
    axios
      .post(`${baseUrl}create-new-users`, payload, {
        headers: {
          "Content-Type": "application/json",
        },
      })
      .then((response) => {
        console.log(response);
        if (response.data.status === "success") {
          navigate("/", { replace: true });
        } else {
          console.log(
            "Error Occurred while creating new user: ",
            response.data
          );
        }
      })
      .catch((err) =>
        console.log("Error occurred while creating users: ", err)
      );
  };
  return (
    <div className={style.Login}>
      <form className={style.RegisterCentre}>
        <span className={style.RegisterSpan}>BMS App</span>
        <input
          type="text"
          id="Name"
          placeholder="Name"
          autoComplete="off"
          onChange={(e) => setName(e.target.value)}
          className={style.RegisterInput}
        />
        <input
          type="text"
          id="UserID"
          placeholder="UserID"
          onChange={(e) => setUserId(e.target.value)}
          className={style.RegisterInput}
        />
        <input
          type="email"
          id="Email"
          placeholder="Email"
          autoComplete="offs"
          onChange={(e) => setEmail(e.target.value)}
          className={style.RegisterInput}
        />
        <input
          type="tel"
          id="ContactNo"
          placeholder="Contact No"
          onChange={(e) => setContactNo(e.target.value)}
          className={style.RegisterInput}
        />
        <input
          type={visible ? "text" : "password"}
          id="Password"
          placeholder="Password"
          autoComplete="off"
          onChange={(e) => setPassword(e.target.value)}
          className={style.RegisterInput}
        />
        <span className={style.RegisterPwdIcon}>
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
          onClick={(e) => {
            e.preventDefault();
            if (!email.includes("@")) {
              setError("Email doesn't contain @!");
            } else if (contactNo.length !== 10) {
              setError("ContactNo should contain max 10 digits!");
            } else if (contactNo.charAt(0) === "0") {
              setError("ContactNo doesn't start with zero!");
            } else if (isNaN(Number(contactNo))) {
              setError("ContactNo should contain only digits!");
              navigate("/register");
            } else {
              sendRegisterReq({
                Name: name,
                Email: email,
                UserId: userId,
                ContactNo: contactNo,
                Password: password,
                Category: "user",
              });
            }
          }}
        >
          Register
        </button>
        <span>
          Have an account, <Link to="/">Login</Link>
        </span>
      </form>
    </div>
  );
}

export default Register;
