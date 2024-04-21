import React from "react";

import './ErrorMessage.css'

const ErrorMessage: React.FC<{ title: string, message: string }> = ({title, message}) => {
    return (
        <div className="ErrorMessage">
            <p className="ErrorMessage_err-title">{title}</p>
            <label className="ErrorMessage_err-msg">{message}</label>
        </div>
    )
}

export default ErrorMessage