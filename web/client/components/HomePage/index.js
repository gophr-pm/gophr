
import classNames from 'classnames'
import { Link } from 'react-router'
import React, { Component } from 'react'

import style from './style.css'
import NavLinks from '../NavLinks'
import SearchBox from '../SearchBox'

class HomePage extends Component {
  state = {
    searchQuery:                  '',
    heroMessageHighlightsVisible: false
  }

  componentDidMount() {
    setTimeout(() => this.setState({
      heroMessageHighlightsVisible: true,
    }), 500)
  }

  onSearchQueryChanged(searchQuery) {
    this.setState({ searchQuery })
  }

  render() {
    const { searchQuery, heroMessageHighlightsVisible } = this.state

    const heroMessageHighlightClasses = classNames(style.heroMessageHighlight, {
      [style.heroMessageHighlight__visible]: heroMessageHighlightsVisible
    })

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
              <a
                href="https://github.com/gophr-pm/gophr"
                target="_blank"
                className={heroMessageHighlightClasses}>gophr</a>
              <span>&nbsp;is the package manager that&nbsp;</span>
              <a
                href="https://golang.org"
                target="_blank"
                className={heroMessageHighlightClasses}>Go</a>
              <span>&nbsp;deserves.</span>
            </div>
            <div className={style.heroButtons}>
              <div
                style={{ backgroundColor: '#5bbb8d' }}
                className={style.heroButton}>Learn More</div>
              <div
                style={{ backgroundColor: '#4b7a7b' }}
                className={style.heroButton}>Find Packages</div>
            </div>
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
