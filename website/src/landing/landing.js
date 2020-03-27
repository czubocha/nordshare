import React from 'react';
import Button from 'react-bootstrap/Button'
import {Link} from "react-router-dom";

const Landing = () => {
    return (
        <div className="App">
            <h1>Nordshare</h1>
            <p>Share notes securely. Simple. {'\u2728'}</p>
            <p><Link to="/notes"><Button variant="primary">Create some</Button></Link></p>
        </div>
    );
};

export default Landing
