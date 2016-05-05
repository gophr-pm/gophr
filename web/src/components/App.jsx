import React from 'react';
import { Link } from 'react-router'
import AppBar from 'material-ui/lib/app-bar';
import AutoComplete from 'material-ui/lib/auto-complete';


const packages = [
  'EasyAPI',
  'Sockets',
  'NeuralNetwork',
];

export default React.createClass({
  render: function() {
    return <div className="App">
      <AppBar
        title="Gophr"
        iconClassNameRight="muidocs-icon-navigation-expand-more"
      >
      <ul>
        <li><Link to="/about">About</Link></li>
      </ul>
      </AppBar>
      <div>
        <AutoComplete
          floatingLabelText="find GO packages"
          filter={AutoComplete.fuzzyFilter}
          dataSource={packages}
        />
      </div>
      <div>
        {this.props.children}
      </div>
      <div>
        <ul>
          <li><Link to="/about">About</Link></li>
          <li><Link to="/support">Support</Link></li>
          <li><Link to="/tokens">Tokens</Link></li>
          <li><Link to="/tutorial">Tutorial</Link></li>
        </ul>
      </div>
    </div>;
  }
});
