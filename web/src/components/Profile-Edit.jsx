import React from 'react';

export default React.createClass({
  getPair: function() {
    return this.props.pair || [];
  },
  render: function() {
    return <div className="Profile-Edit">
      <h1>Profile Edit</h1>
    </div>;
  }
});
