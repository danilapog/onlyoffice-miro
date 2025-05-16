import '@app/global.css';

const { board } = window.miro;

async function init() {
  board.ui.on('icon:click', async () => {
    await board.ui.openPanel({ url: 'application.html' });
  });
}

init();
