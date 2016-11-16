
import React, { Component } from 'react'

import style from './style.css'

const NavLinks = props => (
  <div className={style.main}>
    <a href="https://google.com" className={style.link}>Docs</a>
    <a href="https://google.com" className={style.link}>About</a>
    <a href="https://google.com" className={style.link}>Install</a>
    <a href="https://google.com" className={style.link}>Contribute</a>
  </div>
)

export default NavLinks
