import React from "react";

export default function ImageHolder(props) {
  return (
    <div className="ImageHolder">
      {props.pname}
      {props.username}
    </div>
  );
}
