import { Like } from './like.model';

export class LikeService {

  public likeArray = [
    new Like(1, 1, 1),
    new Like(2, 1, 2),
    new Like(3, 1, 5),
    new Like(1, 2, 5)
  ];
  public getLikeByArtID(id: number): Like[] {
    /* get from API */

    return this.likeArray.filter( el => {
      return el.idA === id;
    });
  }


}

