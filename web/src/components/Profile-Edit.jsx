import React from 'react';

export default React.createClass({
  getPair: function() {
    return this.props.pair || [];
  },
  render: function() {
    return <div className="Profile-Edit">
      <h1>Profile Edit</h1>

      <label>full name</label>
      <input type="text" id="fullname" name="fullname" / >
      <p class="help-text example"></p>

      <label>homepage</label>
      <input type="text" id="homepage" name="homepage" / >
      <p class="help-text example"></p>

      <label>github</label>
      <input type="text" id="github" name="github" / >
      <p class="help-text example"></p>

      <label>twitter</label>
      <input type="text" id="twitter" name="twitter" / >
      <p class="help-text example"></p>

    </div>;
  }
});
