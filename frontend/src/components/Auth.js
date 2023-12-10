import React, { useEffect, useState } from "react";
import { Link,useNavigate } from "react-router-dom";
import axios from "axios"
export default function Auth(props) {
  const navigate = useNavigate()
  const [image, setImage] = useState(
    "https://static-gcp.freepikcompany.com/web-app/media/wepik-1-2000.webp"
  );
  const [place, setPlace] = useState("/Login");
  const [error, setError] = useState("");
  useEffect(() => {
    if ((props.title || " ").toLowerCase() === "signup") {
      setImage(
        "https://static-gcp.freepikcompany.com/web-app/media/wepik-2-2000.webp"
      );
      setPlace("/Login");
    } else {
      setImage(
        "https://static-gcp.freepikcompany.com/web-app/media/wepik-1-2000.webp"
      );
      setPlace("/Signup");
    }
  }, [props.title]);
  const [email,setEmail] = useState("")
  const [pass,setPass] = useState("")
  const handle =async ()=>{
    if(props.title.toLowerCase()==='login'){
      const postData = {
        email: email,
        password: pass,
      };
      const res = await axios.post("http://localhost:8080/api/v1/login",postData);
      const jwt = res.data.token
      localStorage.setItem("Auth",jwt);
      navigate("/")
    }
    if(props.title.toLowerCase()==="signup"){
      const postData = {
        email: email,
        password: pass,
      };
      const res = await axios.post("http://localhost:8080/api/v1/signup",postData);
      if(res.data.success){
        navigate("/")
      }
      else{
        setError(res.data.error)
      }
    }
  }
  return (
    <>
      <Link className="sign" to={place}>
        Go to {place.slice(1, place.length)}
      </Link>
      <div className="form">
        <div className="inner2">
          <h1 className="heading">{props.title}</h1>
        </div>
        <div className="inner">
          <input type="text" name="email" placeholder="Email"value={email} onChange={(e)=>{setEmail(e.target.value)}} />
          <input type="password" name="pass" placeholder="password" value={pass} onChange={(e)=>{setPass(e.target.value)}} />
          <span className="error-msg">{error}</span>
          <input type="submit" value="Submit" onClick={handle}/>
        </div>
      </div>
      <div className="img">
        <img src={image} alt="promo-img" />
      </div>
    </>
  );
}

Auth.defaultProps = {
  title: "login",
};
