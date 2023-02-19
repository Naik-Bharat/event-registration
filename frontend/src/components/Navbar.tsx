function Navbar() {
  return (
    <>
      <nav className="nav">
        <a href="/">
          <img src="logo192.png" />
        </a>
        <ul>
          <li>
            <a href="/auth/google/login">Login</a>
          </li>
        </ul>
      </nav>
    </>
  )
}

export default Navbar;