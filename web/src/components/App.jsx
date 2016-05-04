import React from 'react';

export default React.createClass({
  render: function() {
    return <div className="App">
      <h1>App</h1>
      <ul>
        <li><Link to="/about">About</Link></li>
      </ul>
    </div>;
  }
});
