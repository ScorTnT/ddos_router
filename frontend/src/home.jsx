import React, { useState } from 'react';
import { BrowserRouter as Router, Route, Navigate, Routes } from 'react-router-dom';

function App() {
    const [routing, setRouting] = userState(['networkSetting', 'userSetting', 'intranetSetting']);

    return (
        <>
            <Route>

                <div className="list">
                    <Link to="/network">{ routing[0] }</Link>
                    <hr/>
                </div>

                <div className="list">
                    <Link to="/user">{ routing[1] }</Link>
                    <hr/>
                </div>

                <div className="list">
                    <Link to="/intranet">{ routing[2] }</Link>
                    <hr/>
                </div>

                <Routes>
                    <Route path="/network" element={ <NewtworkSetting /> } />
                    <Route path="/user" element={ <UserSettin /> } />
                    <Route path="/intranet" element={ <NewtworkSetting /> } />
                </Routes>
            </Route>
        </>
    )
}