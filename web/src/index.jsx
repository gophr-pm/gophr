import React from 'react';
import ReactDOM from 'react-dom';
import {RaisedButton} from 'material-ui';
//import injectTapEventPlugin from 'react-tap-event-plugin';

import Voting from './components/Home';

ReactDOM.render((
  <Router>
    <Route path="/" component={Home} />
  </Router>
)
  document.getElementById('app')
);
