import React from 'react';

export default React.createClass({
  getPair: function() {
    return this.props.pair || [];
  },
  render: function() {
    return <div className="Password-Edit">
      <h1>Password Edit</h1>
    </div>;
  }
});
