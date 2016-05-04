import React from 'react';

export default React.createClass({
  getPair: function() {
    return this.props.pair || [];
  },
  render: function() {
    return <div className="Support">
      <h1>Support</h1>
    </div>;
  }
});
