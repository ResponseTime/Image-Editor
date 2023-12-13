import React from "react";

export default function ImageHolder(props) {
  return (
    <div className="ImageHolder">
      <span>Project Name: {props.pname}</span>
      <span>{props.username}</span>
      <img
        src="https://static-gcp.freepikcompany.com/web-app/media/wepik-2-2000.webp"
        alt=""
      />
    </div>
  );
}
