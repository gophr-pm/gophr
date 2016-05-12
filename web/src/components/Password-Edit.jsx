import React from 'react';
import TextField from 'material-ui/lib/TextField';

export default React.createClass({
  getPair: function() {
    return this.props.pair || [];
  },
  render: function() {
    return <div className="Password-Edit">
      <h2>Password Edit</h2>
      <p>An email will be sent to the new address to confirm. Another email will also be sent to the old yasabere@gmail.com address to verify that this is intentional.</p>
      <p>
        <label>Current password</label>
        <br/>
        <TextField
          type="text"
          id="CurrentPassword"
          name="currentpassword"
          hintText=""
          />
      </p>

      <p>
        <label>New password</label>
        <br/>
        <TextField
          type="text"
          id="password"
          name="password"
          hintText=""
          />
      </p>

      <p>
        <label>Verify password</label>
        <br/>
        <TextField
          type="text"
          id="password2"
          name="password2"
          hintText=""
          />
      </p>
    </div>;
  }
});
