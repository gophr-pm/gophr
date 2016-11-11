
import classNames from 'classnames'
import { Link } from 'react-router'
import React, { Component } from 'react'

import style from './style.css'
import NavLinks from '../NavLinks'
import SearchBox from '../SearchBox'

class HomePage extends Component {
  state = {
    searchQuery: ''
  }

  onSearchQueryChanged(searchQuery) {
    this.setState({ searchQuery })
  }

  render() {
    const { searchQuery } = this.state

    return (
      <div className={style.main}>
        <div className={style.header}>
          <div className={style.left}>
            <Link to="/" className={style.logo}></Link>
            <SearchBox
              query={searchQuery}
              onQueryChanged={::this.onSearchQueryChanged} />
          </div>
          <div className={style.right}>
            <NavLinks />
          </div>
        </div>
        <div className={style.splash}>
          <div className={style.hero}>
            <div className={style.heroMessage}>
              gophr is for gophers.
            </div>
            <SearchBox
              query={searchQuery}
              onQueryChanged={::this.onSearchQueryChanged} />
          </div>
        </div>
        <div className={style.section}>
          section
        </div>
      </div>
    )
  }
}

export default HomePage
