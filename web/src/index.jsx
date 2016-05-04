import React from 'react';
import ReactDOM from 'react-dom';
import {RaisedButton} from 'material-ui';
import { Router, Route, Link } from 'react-router'
//import injectTapEventPlugin from 'react-tap-event-plugin';

import About from './components/About';
import EmailEdit from './components/Email-Edit';
import Home from './components/Home';
import Package from './components/Package';
import PasswordEdit from './components/Password-Edit';
import ProfileEdit from './components/Profile-Edit';
import Profile from './components/Profile';
import SubscriptionsEdit from './components/Subscriptions-Edit';
import Support from './components/Support';
import Tokens from './components/Tokens';
import Tutorial from './components/Tutorial';
import NoMatch from './components/Home';

ReactDOM.render((
  <Router>
    <Route path="/" component={Home} >
      <Route path="About" component={About} />
      <Route path="Email-Edit" component={EmailEdit} />
      <Route path="Home" component={Home} />
      <Route path="Package" component={Package} >
        <Route path="/:name" />
      </Route>
      <Route path="Password-Edit" component={PasswordEdit} />
      <Route path="Profile" component={Profile} />
      <Route path="Subscriptions-Edit" component={SubscriptionsEdit} />
      <Route path="Support" component={Support} />
      <Route path="Tokens" component={Tokens} />
      <Route path="Tutorial" component={Tutorial} />
    </Route>
  </Router>
  ),document.getElementById('app')
);
