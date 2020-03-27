import React from 'react';
import './App.css';
import Landing from './landing/landing'
import Note from "./note/note";
import {BrowserRouter as Router, Route, Switch} from "react-router-dom";
import Jumbotron from "react-bootstrap/Jumbotron";

function App() {
    return (
            <header className="App-header">
                <Router>
                    <Jumbotron>
                        <Switch>
                            <Route path="/notes" component={Note}/>
                            <Route path="/" component={Landing}/>
                        </Switch>
                    </Jumbotron>
                </Router>
            </header>
    );
}

export default App;
