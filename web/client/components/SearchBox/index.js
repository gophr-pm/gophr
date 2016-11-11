
import classNames from 'classnames'
import React, { PropTypes, Component } from 'react'

import style from './style.css'

class SearchBox extends Component {
  static propTypes = {
    query:          PropTypes.string,
    onQueryChanged: PropTypes.func,
  }

  input = null
  state = {
    focused: false
  }

  handleChromeClick() {
    if (this.input) {
      this.input.focus()
    }
  }

  handleChange(e) {
    if (this.props.onQueryChanged) {
      this.props.onQueryChanged(e.target.value)
    }
  }

  handleFocus() {
    this.setState({ focused: true })
  }

  handleBlur() {
    this.setState({ focused: false })
  }

  render() {
    const { query } = this.props
    const { focused } = this.state

    const hintClasses = classNames(style.hint, {
      [style.hint__hidden]: focused || query
    })
    const iconImageClasses = classNames(style.iconImage, {
      [style.iconImage__highlighted]: focused
    })

    return (
      <div className={style.main} onClick={::this.handleChromeClick}>
        <input
          ref={r => this.input = r}
          type="text"
          value={query}
          onBlur={::this.handleBlur}
          onFocus={::this.handleFocus}
          onChange={::this.handleChange}
          className={style.input} />
        <div className={hintClasses}>Search for packages</div>
        <div className={style.icon}>
          <svg className={iconImageClasses} viewBox="0 0 24 24">
            <path d="M15.5 14h-.79l-.28-.27C15.41 12.59 16 11.11 16 9.5 16 5.91 13.09 3 9.5 3S3 5.91 3 9.5 5.91 16 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 4.99L20.49 19l-4.99-5zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14z" />
            <path d="M0 0h24v24H0z" fill="none" />
          </svg>
        </div>
      </div>
    )
  }
}

export default SearchBox
