import React from 'react';
import ReactDOM from 'react-dom';
import {RaisedButton} from 'material-ui';
import { Router, Route, Link, browserHistory  } from 'react-router'
import {createStore, applyMiddleware} from 'redux';
import {Provider} from 'react-redux';
import reducer from './reducer';
import { getProfile } from './actions';
import { callAPIMiddleware } from './middleware/API';

//import injectTapEventPlugin from 'react-tap-event-plugin';

//import css
//import './mainold.html'
import './stylesheets/styles.css';
//console.log('__dirname', __dirname + 'css/styles.scss')

//import app components
import App from './components/App';
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


const createStoreWithMiddleware = applyMiddleware(
  callAPIMiddleware
)(createStore);
const store = createStoreWithMiddleware(reducer);
//store.dispatch(setClientId(getClientId()));
store.dispatch(getProfile("yasabere"));

const routes = <Route path="" component={App} >
      <Route path="About" component={About} />
      <Route path="Email-Edit" component={EmailEdit} />
      <Route path="Home" component={Home} />
      <Route path="Package" component={Package} >
        <Route path=":name" component={Package} />
      </Route>
      <Route path="Password-Edit" component={PasswordEdit} />
      <Route path="Profile" component={Profile} />
      <Route path="Profile-Edit" component={ProfileEdit} />
      <Route path="Subscriptions-Edit" component={SubscriptionsEdit} />
      <Route path="Support" component={Support} />
      <Route path="Tokens" component={Tokens} />
      <Route path="Tutorial" component={Tutorial} />
      <Route path="/" component={Home} />
    </Route>;

ReactDOM.render(
  <Provider store={store}>
    <Router >{routes}</Router>
  </Provider>,
  document.getElementById('app')
);

//history={hashHistory}
