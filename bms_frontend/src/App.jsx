import Dashboard from "./Pages/Dashboard";
import { Navigate, Route, Routes, useLocation } from "react-router-dom";
import LendCart from "./Pages/LendCart";
import ReturnCart from "./Pages/ReturnCart";
import BorrowHistory from "./Pages/BorrowHistory";
import Account from "./Pages/Account";
import Notifications from "./Pages/Notifications";
import NoPage from "./Pages/NoPage";
import Login from "./Pages/Login";
import Register from "./Pages/Register";
import Genre from "./Pages/Genre";

function App() {
  let location = useLocation();
  const isLoggedIn = location.state && location.state.isLoggedIn;
  const userInfo = location.state && location.state.userInfo;
  return (
    <div className="App">
      <Routes>
        <Route element={<Login />} path="/"></Route>
        <Route element={<Register />} path="/register"></Route>
        <Route
          element={
            isLoggedIn ? (
              <Dashboard isLoggedIn={isLoggedIn} userInfo={userInfo} />
            ) : (
              <Navigate to="/" replace="true" />
            )
          }
          path="/dashboard"
        ></Route>
        <Route
          element={
            isLoggedIn ? (
              <LendCart isLoggedIn={isLoggedIn} userInfo={userInfo} />
            ) : (
              <Navigate to="/" replace="true" />
            )
          }
          path="/lendCart"
        ></Route>
        <Route
          element={
            isLoggedIn ? (
              <ReturnCart isLoggedIn={isLoggedIn} userInfo={userInfo} />
            ) : (
              <Navigate to="/" replace="true" />
            )
          }
          path="/returnCart"
        ></Route>
        <Route
          element={
            isLoggedIn ? (
              <BorrowHistory isLoggedIn={isLoggedIn} userInfo={userInfo} />
            ) : (
              <Navigate to="/" replace="true" />
            )
          }
          path="/borrowHistory"
        ></Route>
        <Route
          element={
            isLoggedIn ? (
              <Account isLoggedIn={isLoggedIn} userInfo={userInfo} />
            ) : (
              <Navigate to="/" replace="true" />
            )
          }
          path="/account"
        ></Route>
        <Route
          element={
            isLoggedIn ? (
              <Notifications isLoggedIn={isLoggedIn} userInfo={userInfo} />
            ) : (
              <Navigate to="/" replace="true" />
            )
          }
          path="/notifications"
        ></Route>
        <Route
          path="/genre/:id"
          element={
            isLoggedIn ? (
              <Genre isLoggedIn={isLoggedIn} userInfo={userInfo} />
            ) : (
              <Navigate to="/" replace="true" />
            )
          }
        ></Route>
        <Route path="*" element={<NoPage />}></Route>
      </Routes>
    </div>
  );
}

export default App;
