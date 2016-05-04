import React from 'react';

export default React.createClass({
  getPair: function() {
    return this.props.pair || [];
  },
  render: function() {
    return <div className="Email-Edit">
      <h1>Email Edit</h1>
    </div>;
  }
});
