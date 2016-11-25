
import classNames from 'classnames'
import { Link } from 'react-router'
import React, { Component } from 'react'

import style from './style.css'

class HomePage extends Component {
  state = {
    heroMessageHighlightsVisible: false
  }

  componentDidMount() {
    setTimeout(() => this.setState({
      heroMessageHighlightsVisible: true,
    }), 500)
  }

  render() {
    const { heroMessageHighlightsVisible } = this.state

    const heroMessageHighlightClasses = classNames(style.heroMessageHighlight, {
      [style.heroMessageHighlight__visible]: heroMessageHighlightsVisible
    })

    return (
      <div className={style.main}>
        <div className={style.siteComingSoon}>Site Coming Soon</div>
        <div className={style.splash}>
          <div className={style.logo}>
            <div className={style.logoLabel}>alpha</div>
          </div>
          <div className={style.hero}>
            <div className={style.heroMessage}>
              <a
                href="https://github.com/gophr-pm/gophr"
                target="_blank"
                className={heroMessageHighlightClasses}>gophr</a>
              <span> is the package manager that </span>
              <a
                href="https://golang.org"
                target="_blank"
                className={heroMessageHighlightClasses}>Go</a>
              <span> deserves.</span>
            </div>
            <div className={style.heroButtons}>
              <a
                href="https://docs.gophr.pm"
                style={{ backgroundColor: '#5bbb8d' }}
                className={style.heroButton}>Docs</a>
              <a
                href="https://github.com/gophr-pm/gophr"
                style={{ backgroundColor: '#4b7a7b' }}
                className={style.heroButton}>Repo</a>
            </div>
          </div>
        </div>
      </div>
    )
  }
}

export default HomePage
