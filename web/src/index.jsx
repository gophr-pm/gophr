import React from 'react';
import ReactDOM from 'react-dom';
import {RaisedButton} from 'material-ui';
//import injectTapEventPlugin from 'react-tap-event-plugin';

//import Voting from './components/Voting';

const pair = ['Trainspotting', '28 Days Later'];

var Main = React.createClass({
  getPair: function() {
    return this.props.pair || [];
  },
  render: function() {
    return <div className="voting">
      {this.getPair().map(entry =>
        <button key={entry}>
          <h1>{entry}</h1>
        </button>
      )}
      <RaisedButton label="Default" />
    </div>;
  }
});

ReactDOM.render(
  <Main pair={pair} />,
  document.getElementById('app')
);
