import React from 'react';
import {Route, Switch, useRouteMatch} from "react-router-dom";
import Create from "./create/create";
import Show from "./readModifyDelete/readModifyDelete";

const Note = () => {
    let match = useRouteMatch();
    return (
        <Switch>
            <Route path={`${match.path}/:id`} component={Show}/>
            <Route path={match.path} component={Create}/>
        </Switch>
    )
};

export default Note

